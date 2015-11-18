package main

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"os"

	"github.com/codegangsta/cli"
	"github.com/fatih/color"
	"github.com/nadnerb/cli_command"
	"github.com/nadnerb/terraform_config"
	"github.com/nadnerb/terraform_exec/file"
	"github.com/nadnerb/terraform_exec/security"
	"github.com/nadnerb/terraform_exec/sync"
)

var cyan = color.New(color.FgCyan).SprintFunc()
var bold = color.New(color.FgWhite, color.Bold).SprintFunc()
var green = color.New(color.FgGreen).SprintFunc()
var red = color.New(color.FgRed).SprintFunc()

func main() {

	app := cli.NewApp()
	app.Name = "Terraform exec"
	app.Usage = "Execute terraform commands across environments maintaining state in s3\nDefault project layout\n./<project>\n./config/<environment>"

	app.Flags = []cli.Flag{
		cli.BoolFlag{
			Name:  "verbose",
			Usage: "Show more output",
		},
		cli.BoolFlag{
			Name:  "e, env",
			Usage: "Load AWS credentials from environment",
		},
	}

	// Commands
	app.Commands = []cli.Command{
		{
			Name: "run",
			Flags: []cli.Flag{
				cli.BoolFlag{
					Name:  "no-sync",
					Usage: "Don't perform initial s3 sync",
				},
				cli.StringFlag{
					Name: "security",
					Usage: "security provider, current options <default>, <aws-internal>",
				},
				cli.StringFlag{
					Name: "security-role",
					Usage: "security iam role if using -security=aws-internal security",
				},
				cli.StringFlag{
					Name: "config-location",
					Usage: "config location, must be format <location>/<environment>.tfvars",
				},
			},
			Usage:  "Run a terraform command (plan|apply|destroy|refresh)",
			Action: CmdRun,
		},
		{
			Name: "download",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "config-location",
					Usage: "config location, must be format <location>/<environment>.tfvars",
				},
			},
			Usage:  "download existing state to s3",
			Action: CmdDownload,
		},
		{
			Name:   "upload",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "config-location",
					Usage: "config location, must be format <location>/<environment>.tfvars",
				},
			},
			Usage:  "upload existing state to s3",
			Action: CmdUpload,
		},
	}
	app.Run(os.Args)
}

var terraformCommands = map[string]TerraformCommand {
    "plan": TerraformCommand{false, ""},
    "apply": TerraformCommand{true, ""},
    "refresh": TerraformCommand{true, ""},
    "destroy": TerraformCommand{true, "-force"},
}

type TerraformCommand struct {
	sync bool
	extraArgs string
}

func CmdRun(c *cli.Context) {
	if len(c.Args()) != 2 {
		fmt.Printf("Incorrect usage\n")
		fmt.Printf("run <terraform command> <environment>\n")
		os.Exit(1)
	}
	IsTerraformProject()

	terraformCommand := c.Args()[0]
	environment := c.Args()[1]

	security.Apply(c.String("security"), c)
	IsSupportedTerraformCommand(terraformCommand)

	fmt.Println()
	fmt.Println("Run Terraform command")
	fmt.Println("Command:    ", bold(terraformCommand))
	fmt.Println("Environment:", bold(environment))
	fmt.Println()

	configLocation := c.String("config-location")
	config := terraform_config.LoadConfig(c.String("config-location"), environment)

	if c.Bool("no-sync") {
		red("SKIPPING S3 download of current state")
		if command.InputAffirmative() {
			red("Skipped syncing with S3")
		} else {
			DownloadState(config, environment)
		}
	} else {
		DownloadState(config, environment)
	}

	tfVars := terraform_config.TerraformVars(configLocation, environment)
	tfState := terraform_config.TerraformState(environment)
	terraformActions := terraformCommands[terraformCommand]

	// It would be great to use golang terraform so we don't have to install it separately
	// I think we would need to use "github.com/mitchellh/cli" instead of current cli
	fmt.Printf("terraform %s -var-file %s -state=%s -var 'environment=%s' %s\n", terraformCommand, tfVars, tfState, environment, terraformActions.extraArgs)
	fmt.Println("---------------------------------------------")
	fmt.Println()
	cmdName := "terraform"
	cmdArgs := []string{ terraformCommand, "-var-file", tfVars, fmt.Sprintf("-state=%s", tfState) }
	if terraformActions.extraArgs != "" {
		cmdArgs = append(cmdArgs, terraformActions.extraArgs)
	}
	command.Execute(cmdName, cmdArgs)

	fmt.Println("---------------------------------------------")
	if terraformActions.sync {
		fmt.Printf("S3 SYNC new changes\n")
		UploadState(config, environment)
	}
}

func CmdUpload(c *cli.Context) {
	if len(c.Args()) != 1 {
		fmt.Printf("Incorrect usage\n")
		fmt.Printf("upload <environment>\n")
		return
	}
	fmt.Println()
	environment := c.Args()[0]
	config := terraform_config.LoadConfig(c.String("config-location"), environment)

	fmt.Println()
	fmt.Println("Upload Terraform state")
	fmt.Println("Environment:", bold(environment))
	fmt.Println()
	fmt.Printf("Upload current project state to s3\n")

	if command.InputAffirmative() {
		UploadState(config, environment)
	} else {
		color.Red("Aborted")
	}
}

func UploadState(config *terraform_config.AwsConfig, environment string) {
	tfState := terraform_config.TerraformState(environment)
	s3Key := S3Key(config.S3_key, environment)
	fmt.Printf("Uploading project state: %s to: %s/%s\n", green(tfState), green(config.S3_bucket), green(s3Key))

	err := sync.Upload(config.Aws_region, config.S3_bucket, s3Key, tfState)
	fmt.Println()
	if err != nil {
		command.Error("Failed to upload", err)
	} else {
		color.Green("Uploaded successfully to S3")
		fmt.Println()
	}
}

func CmdDownload(c *cli.Context) {
	if len(c.Args()) != 1 {
		fmt.Printf("Incorrect usage\n")
		fmt.Printf("download <environment>\n")
		return
	}
	environment := c.Args()[0]
	fmt.Println()
	fmt.Println("Download Terraform state")
	fmt.Println("Environment:", bold(environment))
	fmt.Println()
	config := terraform_config.LoadConfig(c.String("config-location"), environment)
	DownloadState(config, environment)
}

func DownloadState(config *terraform_config.AwsConfig, environment string) {

	fmt.Println("Syncing project state with S3")

	tfState := terraform_config.TerraformState(environment)
	s3Key := S3Key(config.S3_key, environment)
	fmt.Printf("Downloading project state: %s/%s to: %s\n", cyan(config.S3_bucket), cyan(s3Key), cyan(tfState))

	err := sync.Download(config.Aws_region, config.S3_bucket, s3Key, tfState)
	fmt.Println()
	if err != nil {
		command.Warn("Failed to download", err)
	} else {
		color.Green("Downloaded successfully from S3")
		fmt.Println()
	}
	//sync.Download(&sync.AwsConfig{S3_bucket: config.S3_bucket, S3_key: projectState, Region: config.Region}, fmt.Sprintf("./tfstate/%s/terraform.tfstate")
}

func TextInputAffirmative() string {
	fmt.Print("Are you sure? (yes)\n")
	reader := bufio.NewReader(os.Stdin)
	text, _ := reader.ReadString('\n')
	return text
}

func S3Key(keyName string, environment string) string {
	return fmt.Sprintf("%s/tfstate/%s/terraform.tfstate", keyName, environment)
}

func IsTerraformProject() {
	hasTfFiles, err := file.DirectoryContainsWithExtension(".", ".tf")
	if err != nil && !hasTfFiles {
		command.Error("No terraform files exist", err)
	}
}

func IsSupportedTerraformCommand(terraformCommand string) {
	if _, ok := terraformCommands[terraformCommand]; !ok {
		fmt.Printf("Incorrect usage\n")
		var buffer bytes.Buffer
		buffer.WriteString("Valid commands: ")
		for key := range terraformCommands {
			buffer.WriteString(fmt.Sprintf("%s ", key))
		}
		command.Error("Incorrect run command", errors.New(buffer.String()))
	}
}

package main

import (
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/urfave/cli"
	"github.com/nadnerb/cli_command"
	"github.com/nadnerb/terraform_config"
	"github.com/nadnerb/terraform_exec/security"
	"github.com/nadnerb/terraform_exec/sync"
)

func main() {

	app := cli.NewApp()
	app.Name = ProjectName
	app.Usage = Usage
	app.Version = Version

	app.Flags = []cli.Flag{
		cli.BoolFlag{
			Name:  "verbose",
			Usage: "show more output",
		},
		cli.BoolFlag{
			Name:  "e, env",
			Usage: "load AWS credentials from environment",
		},
		cli.BoolFlag{
			Name:  "commit",
			Usage: "compiled git commit",
		},
	}

	app.Action = func(ctx *cli.Context) error {
		args := ctx.Args()
		if ctx.Bool("commit") {
			CommitMessage()
		} else if args.Present() {
			cli.ShowCommandHelp(ctx, args.First())
		} else {
			cli.ShowAppHelp(ctx)
		}
		return nil
	}

	// Commands
	app.Commands = []cli.Command{
		{
			Name: "plan",
			Flags: []cli.Flag{
				cli.BoolFlag{
					Name:  "destroy",
					Usage: "plan destroy",
				},
				cli.BoolFlag{
					Name:  "no-sync",
					Usage: "Don't perform initial s3 sync",
				},
				cli.StringFlag{
					Name:  "security",
					Usage: "security provider, current options <default>, <aws-internal>",
				},
				cli.StringFlag{
					Name:  "security-role",
					Usage: "security iam role if using -security=aws-internal security",
				},
				cli.StringFlag{
					Name:  "config-location",
					Usage: "config location, must be format <location>/<environment>.tfvars",
				},
			},
			Usage:  "terraform plan",
			Action: CmdPlan,
		},
		{
			Name: "apply",
			Flags: []cli.Flag{
				cli.BoolFlag{
					Name:  "no-sync",
					Usage: "Don't perform initial s3 sync",
				},
				cli.StringFlag{
					Name:  "security",
					Usage: "security provider, current options <default>, <aws-internal>",
				},
				cli.StringFlag{
					Name:  "security-role",
					Usage: "security iam role if using -security=aws-internal security",
				},
				cli.StringFlag{
					Name:  "config-location",
					Usage: "config location, must be format <location>/<environment>.tfvars",
				},
			},
			Usage:  "terraform apply",
			Action: CmdApply,
		},
		{
			Name: "destroy",
			Flags: []cli.Flag{
				cli.BoolFlag{
					Name:  "force",
					Usage: "Force destroy",
				},
				cli.BoolFlag{
					Name:  "no-sync",
					Usage: "Don't perform initial s3 sync",
				},
				cli.StringFlag{
					Name:  "security",
					Usage: "security provider, current options <default>, <aws-internal>",
				},
				cli.StringFlag{
					Name:  "security-role",
					Usage: "security iam role if using -security=aws-internal security",
				},
				cli.StringFlag{
					Name:  "config-location",
					Usage: "config location, must be format <location>/<environment>.tfvars",
				},
			},
			Usage:  "terraform destroy",
			Action: CmdDestroy,
		},
		{
			Name: "refresh",
			Flags: []cli.Flag{
				cli.BoolFlag{
					Name:  "no-sync",
					Usage: "Don't perform initial s3 sync",
				},
				cli.StringFlag{
					Name:  "security",
					Usage: "security provider, current options <default>, <aws-internal>",
				},
				cli.StringFlag{
					Name:  "security-role",
					Usage: "security iam role if using -security=aws-internal security",
				},
				cli.StringFlag{
					Name:  "config-location",
					Usage: "config location, must be format <location>/<environment>.tfvars",
				},
			},
			Usage:  "terraform refresh",
			Action: CmdRefresh,
		},
		{
			Name: "taint",
			Flags: []cli.Flag{
				cli.BoolFlag{
					Name:  "no-sync",
					Usage: "Don't perform initial s3 sync",
				},
				cli.StringFlag{
					Name:  "module",
					Usage: "tainted module path",
				},
				cli.StringFlag{
					Name:  "security",
					Usage: "security provider, current options <default>, <aws-internal>",
				},
				cli.StringFlag{
					Name:  "security-role",
					Usage: "security iam role if using -security=aws-internal security",
				},
				cli.StringFlag{
					Name:  "config-location",
					Usage: "config location, must be format <location>/<environment>.tfvars",
				},
			},
			Usage:  "terraform taint",
			Action: CmdTaint,
		},
		{
			Name: "untaint",
			Flags: []cli.Flag{
				cli.BoolFlag{
					Name:  "no-sync",
					Usage: "Don't perform initial s3 sync",
				},
				cli.StringFlag{
					Name:  "module",
					Usage: "tainted module path",
				},
				cli.StringFlag{
					Name:  "security",
					Usage: "security provider, current options <default>, <aws-internal>",
				},
				cli.StringFlag{
					Name:  "security-role",
					Usage: "security iam role if using -security=aws-internal security",
				},
				cli.StringFlag{
					Name:  "config-location",
					Usage: "config location, must be format <location>/<environment>.tfvars",
				},
			},
			Usage:  "terraform untaint",
			Action: CmdUntaint,
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
			Name: "upload",
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

type TerraformOperation struct {
	command        string
	environment    string
	config         *terraform_config.AwsConfig
	configLocation string
	tfVars         string
	tfState        string
	extraArgs      string
	args           []string
}

func CmdPlan(c *cli.Context) error {
	operation := initialize(c, "plan")
	if c.Bool("destroy") {
		operation.extraArgs = "-destroy"
	}
	run(operation)
	return nil
}

func CmdApply(c *cli.Context) error {
	operation := initialize(c, "apply")
	run(operation)
	resync(operation)
	return nil
}

func CmdTaint(c *cli.Context) error {
	operation := initialize(c, "taint")
	operation.extraArgs = strings.Join(c.Args()[1:], " ")
	run(operation)
	resync(operation)
	return nil
}

func CmdUntaint(c *cli.Context) error {
	operation := initialize(c, "untaint")
	operation.extraArgs = strings.Join(c.Args()[1:], " ")
	run(operation)
	resync(operation)
	return nil
}

func CmdRefresh(c *cli.Context) error {
	operation := initialize(c, "refresh")
	run(operation)
	resync(operation)
	return nil
}

func CmdDestroy(c *cli.Context) error {
	operation := initialize(c, "destroy")
	if !c.Bool("force") {
		fmt.Println(command.Red("Are you sure you want to destroy your environment?"))
		fmt.Println()
		if !command.InputAreYouSure() {
			fmt.Println(command.Cyan("Terraform destroy cancelled"))
			fmt.Println()
			os.Exit(0)
		}
	}
	fmt.Println(command.Cyan("Destroying environment"))
	fmt.Println()

	operation.extraArgs = "-force"
	run(operation)
	resync(operation)
	return nil
}

func initialize(c *cli.Context, terraformCommand string) TerraformOperation {
	if len(c.Args()) < 1 {
		fmt.Printf("Incorrect usage\n")
		fmt.Printf("%s <environment>\n", terraformCommand)
		os.Exit(1)
	}

	fmt.Println(c.Args())
	environment := c.Args()[0]

	security.Apply(c.String("security"), c)

	fmt.Println()
	fmt.Println("Execute Terraform command")
	fmt.Println("Command:    ", command.Bold(terraformCommand))
	fmt.Println("Environment:", command.Bold(environment))
	fmt.Println()

	configLocation := c.String("config-location")
	config := terraform_config.LoadConfig(c.String("config-location"), environment)

	getState(c.Bool("no-sync"), config, environment)

	tfVars := terraform_config.TerraformVars(configLocation, environment)
	tfState := terraform_config.TerraformState(environment)

	return TerraformOperation{
		command:        terraformCommand,
		environment:    environment,
		tfVars:         tfVars,
		tfState:        tfState,
		config:         config,
		configLocation: configLocation,
		args:           os.Args[2:],
	}
}

func getState(skip bool, config *terraform_config.AwsConfig, environment string) {
	if skip {
		fmt.Printf("%s\n", command.Red("Warning: skipping S3 download of current state"))
		if command.InputAreYouSure() {
			fmt.Println()
			fmt.Printf(command.Cyan("Skipped syncing with S3\n"))
			fmt.Println()
		} else {
			DownloadState(config, environment)
		}
	} else {
		DownloadState(config, environment)
	}
}

func resync(operation TerraformOperation) {
	fmt.Printf("%s\n", command.Green("S3 SYNC new changes"))
	UploadState(operation.config, operation.environment)
}

func run(operation TerraformOperation) {
	cmdName := "terraform"
	cmdArgs := commandArgs(operation)

	fmt.Println("---------------------------------------------")
	fmt.Println(command.Bold(cmdName), command.Bold(strings.Join(cmdArgs, " ")))
	fmt.Println("---------------------------------------------")
	fmt.Println()
	command.Default().Execute(cmdName, cmdArgs)
	fmt.Println()
	fmt.Println("---------------------------------------------")
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
	fmt.Println("Environment:", command.Bold(environment))
	fmt.Println()
	fmt.Printf("Upload current project state to s3\n")

	if command.InputAreYouSure() {
		UploadState(config, environment)
	} else {
		fmt.Println(command.Red("Aborted"))
	}
}

func UploadState(config *terraform_config.AwsConfig, environment string) {
	tfState := terraform_config.TerraformState(environment)
	s3Key := S3Key(config.S3_key, environment)
	fmt.Printf("Uploading project state: %s to: %s/%s\n", command.Green(tfState), command.Green(config.S3_bucket), command.Green(s3Key))

	err := sync.Upload(config.Aws_region, config.S3_bucket, s3Key, tfState)
	fmt.Println()
	if err != nil {
		command.Error("Failed to upload", err)
	} else {
		fmt.Println(command.Green("Uploaded successfully to S3"))
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
	fmt.Println("Environment:", command.Bold(environment))
	fmt.Println()
	config := terraform_config.LoadConfig(c.String("config-location"), environment)
	DownloadState(config, environment)
}

func DownloadState(config *terraform_config.AwsConfig, environment string) {

	fmt.Println("Syncing project state with S3")

	tfState := terraform_config.TerraformState(environment)
	s3Key := S3Key(config.S3_key, environment)
	fmt.Printf("Downloading project state: %s/%s to: %s\n", command.Cyan(config.S3_bucket), command.Cyan(s3Key), command.Cyan(tfState))

	err := sync.Download(config.Aws_region, config.S3_bucket, s3Key, tfState)
	fmt.Println()
	if err != nil {
		command.Warn("Failed to download", err)
		if !command.InputAffirmative("Continue without using syncronised state?") {
			fmt.Println(command.Cyan("Operation cancelled"))
			os.Exit(0)
		}
	} else {
		fmt.Println(command.Green("Downloaded successfully from S3"))
	}
	fmt.Println()
}

func S3Key(keyName string, environment string) string {
	return fmt.Sprintf("%s/tfstate/%s/terraform.tfstate", keyName, environment)
}

func terraformArgs(arguments []string) string {
	args := filter(arguments, func(s string) bool {
		match, _ := regexp.MatchString("^-([a-z]+).*", s)
		return match
	})
	return strings.Join(args, " ")
}

func filter(s []string, fn func(string) bool) []string {
	var p []string // == nil
	for _, v := range s {
		if fn(v) {
			p = append(p, v)
		}
	}
	return p
}

func commandArgs(operation TerraformOperation) []string {
	cmdArgs := []string{
		operation.command,
		"-var-file", operation.tfVars,
		fmt.Sprintf("-state=%s", operation.tfState),
		"-var", fmt.Sprintf("environment=%s", operation.environment),
		"-var", fmt.Sprintf("config_location=%s", operation.configLocation),
	}
	args := terraformArgs(operation.args)
	if args != "" {
		cmdArgs = append(cmdArgs, args)
	}
	if operation.extraArgs != "" {
		cmdArgs = append(cmdArgs, operation.extraArgs)
	}
	return cmdArgs
}

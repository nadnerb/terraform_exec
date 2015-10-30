package main

import (
	"bufio"
	"fmt"
	"github.com/codegangsta/cli"
	"github.com/fatih/color"
	"github.com/nadnerb/terraform_exec/sync"
	"github.com/nadnerb/terraform_exec/util"
	"os"
	"path/filepath"
)

// REQUEST: curl http://169.254.169.254/latest/meta-data/iam/security-credentials/ROLE

var cyan = color.New(color.FgCyan).SprintFunc()
var green = color.New(color.FgGreen).SprintFunc()
var red = color.New(color.FgRed).SprintFunc()

func main() {

	app := cli.NewApp()
	app.Name = "Terraform exec"
	app.Usage = "Execute terraform commands across environments maintaining state in s3\nExpects layout build/<project>\nExpects layout config/<project>"

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
			Name:   "setup",
			Usage:  "setup initial execute configuration location",
			Action: CmdSetup,
		},
		{
			Name:   "upload",
			Usage:  "upload existing state to s3",
			Action: CmdUpload,
		},
		{
			Name: "run",
			Flags: []cli.Flag{
				cli.BoolFlag{
					Name:  "no-sync",
					Usage: "Don't perform initial s3 sync",
				},
				cli.StringFlag{
					Name: "config",
					//Usage: "config location, must be format <location>/<project>?/<environment>.tfvars",
					Usage: "config location, must be format <location>/<environment>.tfvars",
				},
				cli.StringFlag{
					Name:  "projectLocation",
					Usage: "terraform project location",
				},
			},
			Usage:  "Run a terraform command (plan|apply|destroy|refresh)",
			Action: CmdRun,
		},
	}
	app.Run(os.Args)

}

func CmdRun(c *cli.Context) {
	if len(c.Args()) != 2 {
		fmt.Printf("Incorrect usage\n")
		fmt.Printf("<terraform command> <environment>\n")
		return
	}

	fmt.Println()
	command := c.Args()[0]
	environment := c.Args()[1]

	e := c.Bool("e")
	if e {
		fmt.Printf("*****************************")
		fmt.Printf("Using environment settings %s", cyan(e))
	}
	env := c.Bool("env")
	if e {
		fmt.Printf("*****************************")
		fmt.Printf("Using environment settings %s", cyan(env))
	}

	config := c.String("config")
	if len(config) > 0 {
		fmt.Printf("Using config location: %s\n", cyan(config))
	} else {
		config := fmt.Sprintf(".%sconfig", filepath.Separator)
		fmt.Printf("Using default config location: %s\n", cyan(config))
	}
	if _, err := os.Stat(config); os.IsNotExist(err) {
		Error("Directory does not exist", err)
	}

	projectLocation := c.String("projectLocation")
	if len(projectLocation) > 0 {
		fmt.Printf("Using project location: %s\n", cyan(projectLocation))
		os.Chdir(projectLocation)
	} else {
		projectLocation := "."
		fmt.Printf("Using default project location: %s\n", cyan(projectLocation))
	}
	hasTfFiles, err := util.HasFilesWithExtension(projectLocation, "tf")
	if err != nil && !hasTfFiles {
		Error("No terraform files exist", err)
	}

	fmt.Printf("Using terraform config: %s/%s.tfvars\n", config, environment)
	fmt.Println()
	awsConfigValues, err := sync.LoadAwsConfig(fmt.Sprintf("%s/%s.tfvars", config, environment))
	if err != nil {
		Error("Error", err)
	} else {
		fmt.Println("s3 Bucket: ", awsConfigValues.S3_bucket)
	}

	fmt.Println("---------------------------------------------")
	noSync := c.Bool("no-sync")
	if noSync {
		color.Red("SKIPPING s3 download of current state")
		fmt.Print("Are you sure? (yes)")
		reader := bufio.NewReader(os.Stdin)
		text, _ := reader.ReadString('\n')
		if text == "yes" {
			color.Red("Skipped sync")
		} else {

		}
	} else {
		fmt.Printf("S3 sync\n")
		fmt.Sprintf("Syncing tfstate/%s.tfstate", environment)
		fmt.Printf("...\n")
	}
	fmt.Println("---------------------------------------------")

	fmt.Printf("terraform %s -var-file %s/%s.tfvars -state=%s/tfstate/%s/terraform.tfstate\n", command, config, environment, projectLocation, environment)

	fmt.Println("---------------------------------------------")
	fmt.Printf("S3 SYNC new changes\n")
}

// SETUP

func CmdSetup(c *cli.Context) {
	if len(c.Args()) != 1 {
		fmt.Printf("Incorrect usage\n")
		fmt.Printf("<config location>\n")
		return
	}
	f, e := os.Create(".terraform.cfg")
	if e != nil {
		panic(e)
	}
	configLocation := c.Args()[0]
	foo, e := f.WriteString(fmt.Sprintf("%s\n", configLocation))
	if e != nil {
		panic(e)
	}
	fmt.Printf("wrote %d bytes\n", foo)
	f.Sync()
}

func CmdUpload(c *cli.Context) {
	if len(c.Args()) != 1 {
		fmt.Printf("Incorrect usage\n")
		return
	}
	fmt.Printf("Upload current project state to s3")
}

func Error(errorMessage string, error error) {
	fmt.Println()
	fmt.Println("---------------------------------------------")
	fmt.Fprintf(os.Stderr, "ERROR\n")
	fmt.Fprintf(os.Stderr, "%s: %s\n", cyan(errorMessage), red(error))
	fmt.Println("---------------------------------------------")
	fmt.Println()
	os.Exit(1)
}

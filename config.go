package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/nadnerb/terraform_exec/command"
	"github.com/nadnerb/terraform_exec/sync"
)

func LoadConfig(configLocation string, environment string) *sync.AwsConfig {
	tfVars := TerraformVars(configLocation, environment)
	awsConfig, err := sync.LoadAwsConfig(tfVars)
	if err != nil {
		command.Error("Error", err)
	}
	fmt.Printf("Using terraform config: %s\n", cyan(tfVars))
	fmt.Println()
	fmt.Println("AWS credentials")
	fmt.Println("s3 bucket: ", bold(awsConfig.S3_bucket))
	fmt.Println("s3 key:    ", bold(awsConfig.S3_key))
	fmt.Println("aws region:", bold(awsConfig.Aws_region))
	fmt.Println()
	return awsConfig
}

func Config(config string) string {
	if len(config) > 0 {
		if _, err := os.Stat(config); os.IsNotExist(err) {
			command.Error("Directory does not exist", err)
		}
		return config
	} else {
		return DefaultConfig()
	}
}

func DefaultConfig() string {
	defaultConfig, _ := filepath.Abs("./config/")
	fmt.Printf("Using default config location: %s\n", cyan(defaultConfig))
	if _, err := os.Stat(defaultConfig); os.IsNotExist(err) {
		command.Error("Directory does not exist", err)
	}
	return defaultConfig
}


package sync

import (
	"fmt"
	"os"
	"io/ioutil"

	"github.com/hashicorp/hcl"
)

func LoadAwsConfig(path string) (*AwsConfig, error) {
	var value AwsConfig

	if _, err := os.Stat(path); err != nil {
		return nil, err
	}
	err := hcl.Decode(&value, ReadFile(path))
	if err != nil {
		return nil, err
	}
	return &value, nil
}

func ReadFile(path string) string {
	d, err := ioutil.ReadFile(path)
	if err != nil {
		fmt.Errorf(
			"Error parsing %s: %s", path, err)
	}

	return string(d)
}

type AwsConfig struct {
	// s3 bucket name
	S3_bucket string
	// s3 key
	S3_key string
	// aws region
	Aws_region string
	// ssh key
	Key_path string
}


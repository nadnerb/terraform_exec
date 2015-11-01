package sync

import (
	"fmt"
	"io/ioutil"

	"github.com/hashicorp/hcl"
)

func LoadAwsConfig(path string) (*AwsConfig, error) {
	var value AwsConfig

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
	S3_bucket string
	S3_key string
	Region string
	Key_path string
	Filename string
}


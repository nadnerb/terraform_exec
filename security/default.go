package security

import (
	"errors"
	"os"

	"github.com/urfave/cli"
)

type DefaultProvider struct {
}

func (d *DefaultProvider) Apply(c *cli.Context) error {
	if len(os.Getenv("AWS_ACCESS_KEY_ID")) == 0 && len(os.Getenv("AWS_SECRET_ACCESS_KEY")) == 0 {
		return errors.New("either AWS_ACCESS_KEY_ID and/or AWS_SECRET_ACCESS_KEY environment variables not set")
	}
	return nil
}

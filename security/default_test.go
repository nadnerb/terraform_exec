package security

import (
	"flag"
	"os"
	"testing"

	"github.com/urfave/cli"
	"github.com/stretchr/testify/assert"
)

func TestAppliesIfEnvironmentVariablesSet(t *testing.T) {
	os.Setenv("AWS_ACCESS_KEY_ID", "a value")
	os.Setenv("AWS_SECRET_ACCESS_KEY", "a value")

	defaultProvider := DefaultProvider{}
	err := defaultProvider.Apply(DefaultContext())

	assert.Nil(t, err)
}

func TestDoesNotApplyIfEnvironmentVariablesSet(t *testing.T) {
	os.Setenv("AWS_ACCESS_KEY_ID", "")
	os.Setenv("AWS_SECRET_ACCESS_KEY", "")

	defaultProvider := DefaultProvider{}
	err := defaultProvider.Apply(DefaultContext())

	assert.NotNil(t, err)
	assert.Equal(t, "either AWS_ACCESS_KEY_ID and/or AWS_SECRET_ACCESS_KEY environment variables not set", err.Error())
}

func DefaultContext() *cli.Context {
	app := cli.NewApp()
	set := flag.NewFlagSet("test", flag.ContinueOnError)
	context := cli.NewContext(app, set, nil)
	return context
}

package security

import (
	"flag"
	"testing"

	"github.com/urfave/cli"
	"github.com/stretchr/testify/assert"
)

// TODO when I figure out how to stub
func TestApplyDefaultSecurity(t *testing.T) {
	//a, _ := Apply("aws-internal", SecurityCliContext())
	assert.True(t, true)
}

func SecurityCliContext() *cli.Context {
	app := cli.NewApp()
	set := flag.NewFlagSet("test", flag.ContinueOnError)
	set.String("security", "aws-internal", "apply security")
	set.String("aws-role", "test-role", "apply security")
	context := cli.NewContext(app, set, nil)
	return context
}


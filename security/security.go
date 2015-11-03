package security

import (
	"github.com/codegangsta/cli"
)

var providers map[string]ProviderFactory

type Provider interface {
	Apply(c *cli.Context) error
}

type ProviderFactory func() (Provider, error)

func init() {
	providers = map[string]ProviderFactory{
		"aws-internal": func() (Provider, error) {
			return &AwsInternalProvider{}, nil
		},
		"default": func() (Provider, error) {
			return &DefaultProvider{}, nil
		},
	}
}

func Apply(providerType string, c *cli.Context) error {
	provider := providers[providerType]
	if provider == nil {
		provider = providers["default"]
	}
	p, err := provider()
	if err != nil {
		return err
	}

	return p.Apply(c)
}

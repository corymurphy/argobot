package argocd

import (
	"github.com/corymurphy/argobot/pkg/cmd"
	"github.com/corymurphy/argobot/pkg/env"
)

type CliClient struct {
	Server        string
	Command       cmd.Command
	ArgoCDBinPath string
	ArgoCliConfig env.ArgoCliConfig
}

func NewCliClient(config env.ArgoCliConfig) *CliClient {
	return &CliClient{
		Server:        config.Server,
		ArgoCDBinPath: config.Command,
		Command:       cmd.NewShellCommand(),
	}
}

func (c *CliClient) Diff(app string, sha string) (string, error) {

	result, err := c.Command.Run(
		c.ArgoCDBinPath,
		"--server",
		c.Server,
		"app",
		"diff",
		"--server-side-generate",
		c.ArgoCliConfig.AdditionalArgs,
		app,
		"--revision",
		sha,
	)

	if err != nil && result == nil {
		return "", err
	}

	return string(result), nil
}

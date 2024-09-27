package argocd

import (
	"fmt"

	"github.com/corymurphy/argobot/pkg/cmd"
	"github.com/corymurphy/argobot/pkg/env"
)

type CliClient struct {
	Server        string
	Command       cmd.Command
	ArgoCDBinPath string
	ArgoConfig    env.ArgoConfig
}

func NewCliClient(config env.ArgoConfig) *CliClient {
	return &CliClient{
		Server:        config.Server,
		ArgoCDBinPath: config.Command,
		ArgoConfig:    config,
		Command:       cmd.NewNonCachingShell(),
	}
}

func (c *CliClient) Plan(app string, sha string) (string, error) {
	command := fmt.Sprintf(
		"%s --server %s app diff --server-side-generate %s --revision %s",
		c.ArgoCDBinPath,
		c.Server,
		app,
		sha,
	)

	if c.ArgoConfig.AdditionalArgs != "" {
		command = command + " " + c.ArgoConfig.AdditionalArgs
	}

	result, err := c.Command.Run(
		"/bin/sh",
		"-c",
		command,
	)

	if err != nil && result == nil {
		return "empty diff", err
	}

	return string(result), nil
}

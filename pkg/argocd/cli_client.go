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
	ArgoCliConfig env.ArgoCliConfig
}

func NewCliClient(config env.ArgoCliConfig) *CliClient {
	return &CliClient{
		Server:        config.Server,
		ArgoCDBinPath: config.Command,
		ArgoCliConfig: config,
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

	if c.ArgoCliConfig.AdditionalArgs != "" {
		command = command + " " + c.ArgoCliConfig.AdditionalArgs
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

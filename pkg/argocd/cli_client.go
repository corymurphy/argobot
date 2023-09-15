package argocd

import (
	"github.com/corymurphy/argobot/pkg/cmd"
)

type CliClient struct {
	Server        string
	Command       cmd.Command
	ArgoCDBinPath string
}

func NewCliClient() *CliClient {
	return &CliClient{
		Server:        "https://argocd.local",
		ArgoCDBinPath: "./bin/argocd-linux-amd64", // TODO: Figure out how to make this dynamic
		Command:       cmd.NewShellCommand(),
	}
}

func (c *CliClient) Diff(app string, sha string) (string, error) {
	// ./bin/argocd-linux-amd64 --server https://argocd.local app diff --server-side-generate helloworld

	result, err := c.Command.Run(
		c.ArgoCDBinPath,
		"--server",
		c.Server,
		"app",
		"diff",
		"--server-side-generate",
		app,
		"--revision",
		sha,
	)

	if err != nil && result == nil {
		return "", err
	}

	return string(result), nil
}

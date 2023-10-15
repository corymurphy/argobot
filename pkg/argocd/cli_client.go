package argocd

import (
	"github.com/corymurphy/argobot/pkg/cmd"
)

type CliClient struct {
	Server        string
	Command       cmd.Command
	ArgoCDBinPath string
	ArgoCliConfig ArgoCliConfig
}

func NewCliClient() *CliClient {
	return &CliClient{
		Server:        "argocd-server.argocd", // TODO: set a good default
		ArgoCDBinPath: "argocd",               // TODO: Figure out how to make this dynamic
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
		"--plaintext",
		app,
		"--revision",
		sha,
	)

	if err != nil && result == nil {
		return "", err
	}

	return string(result), nil
}

package github

import (
	"fmt"

	"github.com/corymurphy/argobot/pkg/cmd"
	"github.com/corymurphy/argobot/pkg/env"
)

type GitClient struct {
	Config        env.Config
	Token         AccessToken
	Command       cmd.Command
	GitBinaryPath string
	Host          string
}

func NewGitClient(config env.Config, token AccessToken) *GitClient {
	return &GitClient{
		Config:        config,
		Token:         token,
		Command:       cmd.NewShellCommand(),
		GitBinaryPath: "git",
		Host:          "github.com",
	}
}

func (g *GitClient) Clone(org string, repo string, path string) error {
	url := fmt.Sprintf(
		"https://x-access-token:%s@%s/%s/%s.git",
		g.Token.Token,
		g.Host,
		org,
		repo,
	)

	_, err := g.Command.Run(
		g.GitBinaryPath,
		"clone",
		url,
		path,
	)

	if err != nil {
		return err
	}

	// if code != 0 {
	// 	return errors.New(fmt.Sprintf("unable to clone repo %s/%s with error code %d", org, repo, code))
	// }

	return nil
}

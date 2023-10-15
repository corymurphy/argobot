// Copyright 2018 Palantir Technologies, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package events

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"strconv"

	"github.com/corymurphy/argobot/pkg/argocd"
	"github.com/corymurphy/argobot/pkg/env"
	cmgithub "github.com/corymurphy/argobot/pkg/github"
	"github.com/corymurphy/argobot/pkg/parsing"
	"github.com/google/go-github/v53/github"
	"github.com/palantir/go-githubapp/githubapp"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
)

type PRCommentHandler struct {
	githubapp.ClientCreator
	Config *env.Config
	// preamble string
}

func (h *PRCommentHandler) Handles() []string {
	return []string{"issue_comment"}
}

func (h *PRCommentHandler) Handle(ctx context.Context, eventType, deliveryID string, payload []byte) error {
	var event github.IssueCommentEvent
	if err := json.Unmarshal(payload, &event); err != nil {
		return errors.Wrap(err, "failed to parse issue comment event payload")
	}

	if !event.GetIssue().IsPullRequest() {
		zerolog.Ctx(ctx).Debug().Msg("Issue comment event is not for a pull request")
		return nil
	}

	repo := event.GetRepo()
	prNum := event.GetIssue().GetNumber()
	installationID := githubapp.GetInstallationIDFromEvent(&event)

	ctx, logger := githubapp.PreparePRContext(ctx, installationID, repo, event.GetIssue().GetNumber())

	logger.Debug().Msgf("event metadata %d %d", installationID, event.GetIssue().GetNumber())

	logger.Debug().Msgf("Event action is %s", event.GetAction())
	if event.GetAction() != "created" {
		return nil
	}

	client, err := h.NewInstallationClient(installationID)
	if err != nil {
		return err
	}

	repoOwner := repo.GetOwner().GetLogin()
	repoName := repo.GetName()

	pr, _, err := client.PullRequests.Get(ctx, repoOwner, repoName, prNum)
	if err != nil {
		return err
	}

	sha := pr.GetHead().SHA

	logger.Debug().Msgf(*sha)

	parser, err := parsing.NewPRCommentParser(event)

	if err != nil {
		logger.Debug().Err(err)
		return nil
	}

	if parser.IsBot {
		logger.Debug().Msg("Issue comment was created by a bot")
		return nil
	}

	if !parser.IsArgoCommand {
		logger.Debug().Msg("comment was not an argo bot command")
		return nil
	}

	if parser.Command == parsing.Help {
		comment := helpComment()
		prComment := github.IssueComment{
			Body: &comment,
		}

		if _, _, err := client.Issues.CreateComment(ctx, repoOwner, repoName, prNum, &prComment); err != nil {
			logger.Error().Err(err).Msg("Failed to comment on pull request")
		}
	}

	if parser.Command != parsing.Diff {
		logger.Debug().Msg("comment was an argo command but only diff is supported")
		return nil
	}

	// create data directory
	dir := filepath.Join(
		h.Config.AppConfig.DataDirectory,
		repoOwner,
		repoName,
		strconv.Itoa(prNum),
	)

	if err := os.MkdirAll(dir, 0755); err != nil {
		logger.Error().Err(err).Msg("Failed to create directory")
	}

	apiClient := cmgithub.NewApiClient(*h.Config)

	token, _ := apiClient.NewAccessToken("39903519")

	git := cmgithub.NewGitClient(*h.Config, token)

	git.Clone(repoOwner, repoName, dir)

	argocd := argocd.NewCliClient()
	// argocd.ArgoCDBinPath = "argocd"
	// sha := event.

	diff, err := argocd.Diff(parser.Application, *sha)

	if err != nil {
		logger.Error().Err(err)
	}

	response := "```diff\n" +
		diff + "\n" +
		"```"

	prComment := github.IssueComment{
		Body: &response,
	}

	if _, _, err := client.Issues.CreateComment(ctx, repoOwner, repoName, prNum, &prComment); err != nil {
		logger.Error().Err(err).Msg("Failed to comment on pull request")
	}

	return nil
}

func helpComment() string {
	return `
Usage: argo COMMAND [ARGS]...

  Allows you to interact with ArgoCD from a Pull Request.

Commands:
  help 		Shows this message
  diff --application myapp --directory relative/git/directory/to/app
`
}

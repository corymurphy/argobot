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
	"fmt"

	"github.com/argoproj/argo-cd/v2/pkg/apiclient/application"
	"github.com/corymurphy/argobot/pkg/argocd"
	"github.com/corymurphy/argobot/pkg/command"
	"github.com/corymurphy/argobot/pkg/env"
	"github.com/corymurphy/argobot/pkg/logging"
	"github.com/corymurphy/argobot/pkg/models"
	"github.com/corymurphy/argobot/pkg/utils"
	"github.com/google/go-github/v53/github"
	"github.com/palantir/go-githubapp/githubapp"
	"github.com/pkg/errors"
)

const maxCommentLength = 32768

type PRCommentHandler struct {
	githubapp.ClientCreator
	Config      *env.Config
	Log         logging.SimpleLogging
	ArgoClient  argocd.ApplicationsClient
	TestingMode bool
}

func (h *PRCommentHandler) Handles() []string {
	return []string{"issue_comment"}
}

func (h *PRCommentHandler) Handle(ctx context.Context, eventType, deliveryID string, payload []byte) error {
	var event github.IssueCommentEvent
	if err := json.Unmarshal(payload, &event); err != nil {
		return errors.Wrap(err, "failed to parse issue comment event payload")
	}

	comment := NewCommentParser(h.Log).Parse(event)
	if (comment.Ignore || comment.ImmediateResponse) && !comment.HasResponseComment {
		return nil
	}

	repo := event.GetRepo()
	prNum := event.GetIssue().GetNumber()
	repoOwner := repo.GetOwner().GetLogin()
	repoName := repo.GetName()

	installationID := githubapp.GetInstallationIDFromEvent(&event)

	ctx, _ = githubapp.PreparePRContext(ctx, installationID, repo, event.GetIssue().GetNumber())

	client, err := h.NewInstallationClient(installationID)
	if err != nil {
		return err
	}

	pr, _, err := client.PullRequests.Get(ctx, repoOwner, repoName, prNum)
	if err != nil {
		return err
	}

	revision := pr.GetHead().SHA
	h.Log.Debug("comment was for sha %s", *revision)

	request := models.NewPullRequestComment(
		*revision,
		models.PullRequest{
			Number: prNum,
			Name:   repoName,
			Owner:  repoOwner,
		},
	)

	if (comment.Ignore || comment.ImmediateResponse) && comment.HasResponseComment {
		prComment := github.IssueComment{
			Body: &comment.CommentResponse,
		}

		if _, _, err := client.Issues.CreateComment(ctx, repoOwner, repoName, prNum, &prComment); err != nil {
			return err
		}
	}
	// TODO: add reaction like this
	// err := e.VCSClient.ReactToComment(logger, baseRepo, pullNum, commentID, e.EmojiReaction)

	// TODO: apply and plan async and respond to hook immediately

	// TODO: move this to startup and cache?

	if comment.Command.Application == "" {
		h.Log.Info(fmt.Sprintf("app %s", comment.Command.Application))

		appResolver := NewApplicationResolver(client, &h.ArgoClient, h.Log)
		apps, err := appResolver.FindApplicationNames(ctx, comment.Command, request.Pull)
		if err != nil {
			h.Log.Err("unable to resolve app name %w", err)
			return err
		}

		// TODO support multiple apps
		comment.Command.Application = apps[0]
	}

	if comment.Command.Name == command.Plan {
		plan, diff, err := h.Plan(ctx, comment.Command.Application, *revision)
		if err != nil {
			h.Log.Err("unable to plan: %w %s", err, plan)
			return err
		}
		var msg string
		if diff {
			msg = fmt.Sprintf("argocd plan for `%s`\n\n", comment.Command.Application) + "```diff\n" + plan + "\n```"
		} else {
			msg = "no diff detected, current state is up to date with this revision."
			h.Log.Info(plan)
		}

		return h.CreateComment(client, ctx, request.Pull, msg, comment.Command.Name.String())
	}

	if comment.Command.Name == command.Apply {
		go func() {
			apply := NewApplyRunner(client, h.Config, h.Log, &h.ArgoClient)
			response, err := apply.Run(ctx, comment.Command, request)
			if err != nil {
				h.Log.Err("unable to apply %w %s", err)
				return
			}
			msg := fmt.Sprintf("apply result for `%s`\n\n", comment.Command.Application) + "```\n" + response.Message + "\n```"
			h.CreateComment(client, ctx, request.Pull, msg, comment.Command.Name.String())
		}()
		return nil
	}

	return errors.Errorf("unsupported argo command")
}

// TODO: this is just temporary while i build the proof of concept
func (h *PRCommentHandler) Plan(ctx context.Context, name string, revision string) (string, bool, error) {
	var plan string
	var diff bool = false
	// var resources argo.ApplicationManagedResourcesResponse // http://localhost:8081/api/v1/applications/helloworld/managed-resources
	// var app argo.ApplicationApplicationResponse

	var resources *application.ManagedResourcesResponse

	resources, err := h.ArgoClient.ManagedResources(name)

	if err != nil {
		return plan, diff, err
	}

	live, err := h.ArgoClient.Get(name)

	if err != nil {
		return plan, diff, err
	}

	target, err := h.ArgoClient.GetManifest(name, revision)
	if err != nil {
		return plan, diff, err
	}

	settings, err := h.ArgoClient.GetSettings()
	if err != nil {
		return plan, diff, err
	}

	// h.Log.Info(l)

	return argocd.Plan(ctx, &settings, live, resources, target, revision, h.Log)
}

// TODO move this to another module
func (h *PRCommentHandler) CreateComment(client *github.Client, ctx context.Context, pull models.PullRequest, comment string, command string) error {
	h.Log.Debug("Creating comment on GitHub pull request %d", pull.Number)
	var sepStart string

	sepEnd := "\n```\n</details>" +
		"\n<br>\n\n**Warning**: Output length greater than max comment size. Continued in next comment."

	if command != "" {
		sepStart = fmt.Sprintf("Continued %s output from previous comment.\n<details><summary>Show Output</summary>\n\n", command) +
			"```diff\n"
	} else {
		sepStart = "Continued from previous comment.\n<details><summary>Show Output</summary>\n\n" +
			"```diff\n"
	}

	truncationHeader := "\n```\n</details>" +
		"\n<br>\n\n**Warning**: Command output is larger than the maximum number of comments per command. Output truncated.\n\n[..]\n"

	comments := utils.SplitComment(comment, maxCommentLength, sepEnd, sepStart, 100, truncationHeader)
	for i := range comments {
		_, resp, err := client.Issues.CreateComment(ctx, pull.Owner, pull.Name, pull.Number, &github.IssueComment{Body: &comments[i]})
		if resp != nil {
			h.Log.Debug("POST /repos/%v/%v/issues/%d/comments returned: %v", pull.Owner, pull.Name, pull.Number, resp.StatusCode)
		}
		if err != nil {
			return err
		}
	}
	return nil
}

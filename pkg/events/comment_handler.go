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
	"math"
	"strings"

	"github.com/corymurphy/argobot/pkg/argocd"
	"github.com/corymurphy/argobot/pkg/command"
	"github.com/corymurphy/argobot/pkg/env"
	"github.com/corymurphy/argobot/pkg/logging"
	"github.com/google/go-github/v53/github"
	"github.com/palantir/go-githubapp/githubapp"
	"github.com/pkg/errors"
)

const maxCommentLength = 32768

type PRCommentHandler struct {
	githubapp.ClientCreator
	Config      *env.Config
	Log         logging.SimpleLogging
	PlanClient  argocd.PlanClient
	ApplyClient argocd.ApplyClient
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

	sha := pr.GetHead().SHA
	h.Log.Debug("comment was for sha %s", *sha)

	if (comment.Ignore || comment.ImmediateResponse) && comment.HasResponseComment {
		prComment := github.IssueComment{
			Body: &comment.CommentResponse,
		}

		if _, _, err := client.Issues.CreateComment(ctx, repoOwner, repoName, prNum, &prComment); err != nil {
			return err
		}
	}

	// TODO: apply and plan async and respond to hook immediately

	if comment.Command.Name == command.Plan {
		plan, err := h.PlanClient.Plan(comment.Command.Application, *sha)
		if err != nil {
			return err
		}
		h.Log.Info(plan)
		return h.CreateBlockComment(client, ctx, repo, prNum, plan, comment.Command.Name.String(), "```diff")
	}

	if comment.Command.Name == command.Apply {
		// TODO: return success and run apply in background
		// TODO: validate that the PR is in an approved/mergeable state
		apply, err := h.ApplyClient.Apply(comment.Command.Application, *sha)
		if err != nil {
			h.Log.Err("argoclient failed while applying %w", err)
			return err
		}
		h.Log.Info(apply)
		return h.CreateComment(client, ctx, repo, prNum, fmt.Sprintf("applied %s sucessfully", comment.Command.Application))
	}

	return errors.Errorf("unsupported argo command")
}

// TODO
func (h *PRCommentHandler) CreateComment(vsc *github.Client, ctx context.Context, repo *github.Repository, pullNum int, comment string) error {
	h.Log.Debug("POST /repos/%s/%s/issues/%d/comments", repo.GetOwner().GetLogin(), repo.GetName(), pullNum)
	_, _, err := vsc.Issues.CreateComment(ctx, repo.GetOwner().GetLogin(), repo.GetName(), pullNum, &github.IssueComment{Body: &comment})
	return err
}

func (h *PRCommentHandler) CreateBlockComment(vsc *github.Client, ctx context.Context, repo *github.Repository, pullNum int, comment string, command string, blockPrefix string) error {
	var sepStart string

	sepEnd := "\n```\n</details>" +
		"\n<br>\n\n**Warning**: Output length greater than max comment size. Continued in next comment."

	if command != "" {
		sepStart = fmt.Sprintf("Continued %s output from previous comment.\n<details><summary>Show Output</summary>\n\n", command) +
			blockPrefix + "\n"
	} else {
		sepStart = "Continued from previous comment.\n<details><summary>Show Output</summary>\n\n" +
			blockPrefix + "\n"
	}

	comments := SplitComment(comment, maxCommentLength, sepEnd, sepStart, blockPrefix)
	for i := range comments {
		h.Log.Debug("POST /repos/%s/%s/issues/%d/comments", repo.GetOwner().GetLogin(), repo.GetName(), pullNum)
		err := h.CreateComment(vsc, ctx, repo, pullNum, comments[i])
		if err != nil {
			return err
		}
	}
	return nil
}

func SplitComment(comment string, maxSize int, sepEnd string, sepStart string, blockPrefix string) []string {
	if strings.TrimSpace(comment) == "" {
		return []string{"No diff detected, resources are up to date."}
	}

	if len(comment) <= maxSize {
		return []string{blockPrefix + "\n" + comment + "\n```"}
	}

	maxWithSep := maxSize - len(sepEnd) - len(sepStart)
	var comments []string
	numComments := int(math.Ceil(float64(len(comment)) / float64(maxWithSep)))
	for i := 0; i < numComments; i++ {
		upTo := min(len(comment), (i+1)*maxWithSep)
		portion := comment[i*maxWithSep : upTo]
		if i < numComments-1 {
			portion += sepEnd
		}
		if i > 0 {
			portion = sepStart + portion
		}
		comments = append(comments, portion)
	}
	return comments
}

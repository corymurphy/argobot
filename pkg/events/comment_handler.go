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

	"github.com/corymurphy/argobot/pkg/argocd"
	"github.com/corymurphy/argobot/pkg/command"
	"github.com/corymurphy/argobot/pkg/env"
	"github.com/corymurphy/argobot/pkg/logging"
	"github.com/corymurphy/argobot/pkg/models"
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

	sha := pr.GetHead().SHA
	h.Log.Debug("comment was for sha %s", *sha)

	request := models.NewPullRequestComment(
		*sha,
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
		plan, err := h.PlanClient.Plan(comment.Command.Application, *sha)
		if err != nil {
			return err
		}
		msg := fmt.Sprintf("argocd plan for `%s`\n\n", comment.Command.Application) + "```diff\n" + plan + "\n```"
		return h.CreateComment(client, ctx, request.Pull, msg, comment.Command.Name.String())
		// prefix
		// return h.CreateBlockComment(client, ctx, request.Pull, prefix, plan, comment.Command.Name.String(), "```diff")

	}

	if comment.Command.Name == command.Apply {
		go func() {
			apply := NewApplyRunner(client, h.Config, h.Log, h.ApplyClient)
			response, err := apply.Run(ctx, comment.Command, request)
			if err != nil {
				h.Log.Err("unable to apply %w", err)
				return
			}
			msg := fmt.Sprintf("apply result for `%s`\n\n", comment.Command.Application) + "```\n" + response.Message + "\n```"
			h.CreateComment(client, ctx, request.Pull, msg, comment.Command.Name.String())
			// plan = "```diff\n" + plan + "\n```"
			// plan = fmt.Sprintf("===== %s ======\n\n", comment.Command.Application) + plan
			// h.CreateComment(client, ctx, request.Pull, plan, comment.Command.Name.String())

			// h.CreateComment(client, ctx, request.Pull, response.Message)
		}()
		return nil
	}

	return errors.Errorf("unsupported argo command")
}

// TODO move this to a separate module
// func (h *PRCommentHandler) CreateComment(vsc *github.Client, ctx context.Context, pull models.PullRequest, comment string) error {
// 	h.Log.Debug("POST /repos/%s/%s/issues/%d/comments", pull.Owner, pull.Name, pull.Number)
// 	_, _, err := vsc.Issues.CreateComment(ctx, pull.Owner, pull.Name, pull.Number, &github.IssueComment{Body: &comment})
// 	return err
// }

// TODO refactor signature
// func (h *PRCommentHandler) CreateComment(vsc *github.Client, ctx context.Context, pull models.PullRequest, prefix string, comment string, command string, blockPrefix string) error {
// 	var sepStart string

// 	sepEnd := "\n```\n</details>" +
// 		"\n<br>\n\n**Warning**: Output length greater than max comment size. Continued in next comment."

// 	if command != "" {
// 		sepStart = fmt.Sprintf("Continued %s output from previous comment.\n<details><summary>Show Output</summary>\n\n", command) +
// 			"```diff\n"
// 	} else {
// 		sepStart = "Continued from previous comment.\n<details><summary>Show Output</summary>\n\n" +
// 			"```diff\n"
// 	}

// 	truncationHeader := "\n```\n</details>" +
// 		"\n<br>\n\n**Warning**: Command output is larger than the maximum number of comments per command. Output truncated.\n\n[..]\n"

// 	comments := SplitComment(comment, maxCommentLength, sepEnd, sepStart, 100, truncationHeader)
// 	for i := range comments {
// 		// if i == 0 {
// 		// 	comments[i] = comments[i]
// 		// }
// 		h.Log.Debug("POST /repos/%s/%s/issues/%d/comments", pull.Owner, pull.Name, pull.Number)
// 		err := h.CreateComment(vsc, ctx, pull, comments[i])
// 		if err != nil {
// 			return err
// 		}
// 	}
// 	return nil
// }

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

	comments := SplitComment(comment, maxCommentLength, sepEnd, sepStart, 100, truncationHeader)
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

func SplitComment(comment string, maxSize int, sepEnd string, sepStart string, maxCommentsPerCommand int, truncationHeader string) []string {
	if len(comment) <= maxSize {
		return []string{comment}
	}

	// No comment contains both sepEnd and truncationHeader, so we only have to count their max.
	maxWithSep := maxSize - max(len(sepEnd), len(truncationHeader)) - len(sepStart)
	var comments []string
	numPotentialComments := int(math.Ceil(float64(len(comment)) / float64(maxWithSep)))
	var numComments int
	if maxCommentsPerCommand == 0 {
		numComments = numPotentialComments
	} else {
		numComments = min(numPotentialComments, maxCommentsPerCommand)
	}
	isTruncated := numComments < numPotentialComments
	upTo := len(comment)
	for len(comments) < numComments {
		downFrom := max(0, upTo-maxWithSep)
		portion := comment[downFrom:upTo]
		if len(comments)+1 != numComments {
			portion = sepStart + portion
		} else if len(comments)+1 == numComments && isTruncated {
			portion = truncationHeader + portion
		}
		if len(comments) != 0 {
			portion = portion + sepEnd
		}
		comments = append([]string{portion}, comments...)
		upTo = downFrom
	}
	return comments
}

// func SplitComment(prefix string, comment string, maxSize int, sepEnd string, sepStart string, blockPrefix string) []string {
// 	if strings.TrimSpace(comment) == "" {
// 		return []string{"No diff detected, resources are up to date."}
// 	}

// 	if len(comment) <= maxSize {
// 		return []string{comment}
// 	}

// 	// if len(blockPrefix+"\n"+comment+"\n```") <= maxSize {
// 	// 	return []string{blockPrefix + "\n" + comment + "\n```"}
// 	// }

// 	maxWithSep := maxSize - len(sepEnd) - len(sepStart)
// 	var comments []string
// 	numComments := int(math.Ceil(float64(len(comment)) / float64(maxWithSep)))
// 	for i := 0; i < numComments; i++ {
// 		upTo := min(len(comment), (i+1)*maxWithSep)
// 		portion := comment[i*maxWithSep : upTo]
// 		if i < numComments-1 {
// 			portion += sepEnd
// 		}
// 		if i > 0 {
// 			portion = sepStart + portion
// 		}
// 		comments = append(comments, portion)
// 	}
// 	return comments
// }

package events

import (
	"context"
	"fmt"

	"github.com/argoproj/argo-cd/v2/pkg/apiclient/application"
	"github.com/corymurphy/argobot/pkg/argocd"
	"github.com/corymurphy/argobot/pkg/command"
	"github.com/corymurphy/argobot/pkg/env"
	vsc "github.com/corymurphy/argobot/pkg/github"
	"github.com/corymurphy/argobot/pkg/logging"
	"github.com/corymurphy/argobot/pkg/utils"
	"github.com/google/go-github/v53/github"
	"github.com/palantir/go-githubapp/githubapp"
	"github.com/pkg/errors"
)

const maxCommentLength = 32768

type PRCommentHandler struct {
	githubapp.ClientCreator
	Config     *env.Config
	Log        logging.SimpleLogging
	ArgoClient argocd.ApplicationsClient
}

func (h *PRCommentHandler) Handles() []string {
	return []string{"issue_comment", "pull_request"}
}

func (h *PRCommentHandler) Handle(ctx context.Context, eventType string, deliveryID string, payload []byte) error {
	// var event github.IssueCommentEvent
	// event.GetInstallation()
	// var pr github.PullRequestEvent
	event, err := vsc.NewEvent(eventType, payload)
	if err != nil {
		h.Log.Err(err, "unable to parse event metadata")
		return nil
	}

	comment := NewCommentParser(h.Log).Parse(event)
	if (comment.Ignore || comment.ImmediateResponse) && !comment.HasResponseComment {
		return nil
	}

	installationID := githubapp.GetInstallationIDFromEvent(&event)

	client, err := h.NewInstallationClient(installationID)
	if err != nil {
		return err
	}

	pr, _, err := client.PullRequests.Get(ctx, event.Repository.Owner, event.Repository.Name, event.PullRequest.Number)
	if err != nil {
		return err
	}

	// TODO: this should happen in NewEvent
	event.Revision = *pr.GetHead().SHA

	if (comment.Ignore || comment.ImmediateResponse) && comment.HasResponseComment {
		prComment := github.IssueComment{
			Body: &comment.CommentResponse,
		}

		if _, _, err := client.Issues.CreateComment(ctx, event.Repository.Owner, event.Repository.Name, event.PullRequest.Number, &prComment); err != nil {
			return err
		}
	}
	// TODO: add reaction like this
	// err := e.VCSClient.ReactToComment(logger, baseRepo, pullNum, commentID, e.EmojiReaction)

	// TODO: apply and plan async and respond to hook immediately

	// TODO: move this to startup and cache?

	if comment.Command.Application == "" {
		appResolver := NewApplicationResolver(client, &h.ArgoClient, h.Log)
		apps, err := appResolver.FindApplicationNames(ctx, comment.Command, event)
		if err != nil {
			h.Log.Err(err, "unable to resolve app name")
			return nil
		}

		comment.Command.Applications = apps
	}

	if comment.Command.Name == command.Plan {
		for _, app := range comment.Command.Applications {
			var err error = nil

			plan, diff, err := h.Plan(ctx, app, event.Revision)
			if err != nil {
				h.Log.Err(err, fmt.Sprintf("unable to plan: %s", plan))
				return err
			}
			var msg string
			if diff {
				msg = fmt.Sprintf("argocd plan for `%s`\n\n", app) + "```diff\n" + plan + "\n```"
			} else {
				msg = "no diff detected, current state is up to date with this revision."
				h.Log.Info(plan)
			}

			err = h.CreateComment(client, ctx, event, msg, comment.Command.Name.String())
			if err != nil {
				h.Log.Err(err, fmt.Sprintf("error while planning %s", app))
			}
		}
		return nil
	}

	// TODO allow autoapply
	if comment.Command.Name == command.Apply {
		go func() {
			apply := NewApplyRunner(client, h.Config, h.Log, &h.ArgoClient)
			response, err := apply.Run(ctx, comment.Command, event)
			if err != nil {
				h.Log.Err(err, "unable to apply")

				return
			}
			msg := fmt.Sprintf("apply result for `%s`\n\n", comment.Command.Application) + "```\n" + response.Message + "\n```"
			h.CreateComment(client, ctx, event, msg, comment.Command.Name.String())
		}()
		return nil
	}

	return errors.Errorf("unsupported argo command")
}

// TODO: this is just temporary while i build the proof of concept
func (h *PRCommentHandler) Plan(ctx context.Context, name string, revision string) (string, bool, error) {
	var plan string
	var diff bool = false
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

	return argocd.Plan(ctx, &settings, live, resources, target, revision, h.Log)
}

// TODO move this to another module
func (h *PRCommentHandler) CreateComment(client *github.Client, ctx context.Context, event vsc.Event, comment string, command string) error {
	h.Log.Debug("Creating comment on GitHub pull request %d", event.PullRequest.Number)
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
		_, resp, err := client.Issues.CreateComment(ctx, event.Repository.Owner, event.Repository.Name, event.PullRequest.Number, &github.IssueComment{Body: &comments[i]})
		if resp != nil {
			h.Log.Debug("POST /repos/%v/%v/issues/%d/comments returned: %v", event.Repository.Owner, event.Repository.Name, event.PullRequest.Number, resp.StatusCode)
		}
		if err != nil {
			return err
		}
	}
	return nil
}

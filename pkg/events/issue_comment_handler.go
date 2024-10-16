package events

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/corymurphy/argobot/pkg/argocd"
	"github.com/corymurphy/argobot/pkg/command"
	"github.com/corymurphy/argobot/pkg/env"
	vsc "github.com/corymurphy/argobot/pkg/github"
	"github.com/corymurphy/argobot/pkg/logging"
	"github.com/google/go-github/v53/github"
	"github.com/palantir/go-githubapp/githubapp"
	"github.com/pkg/errors"
)

var validIssueCommentActions = []string{"created"}

type IssueCommentHandler struct {
	githubapp.ClientCreator
	Config     *env.Config
	Log        logging.SimpleLogging
	ArgoClient argocd.ApplicationsClient
}

func (h *IssueCommentHandler) Handles() []string {
	return []string{"issue_comment"}
}

func (h *IssueCommentHandler) Handle(ctx context.Context, eventType string, deliveryID string, payload []byte) error {
	var issue github.IssueCommentEvent

	if err := json.Unmarshal(payload, &issue); err != nil {
		h.Log.Err(err, "invalid github event payload")
		return fmt.Errorf("invalid github event payload")
	}

	installationID := githubapp.GetInstallationIDFromEvent(&issue)
	client, err := h.NewInstallationClient(installationID)
	if err != nil {
		return err
	}

	pr, _, err := client.PullRequests.Get(ctx, issue.GetRepo().GetOwner().GetLogin(), issue.GetRepo().GetName(), issue.Issue.GetNumber())
	if err != nil {
		h.Log.Err(err, "unable to get revision from pull request")
		return nil
	}
	event := vsc.InitializeFromIssueComment(issue, *pr.GetHead().SHA)

	comment := NewCommentParser(h.Log).Parse(event)
	if (comment.Ignore || comment.ImmediateResponse) && !comment.HasResponseComment {
		return nil
	}

	if (comment.Ignore || comment.ImmediateResponse) && comment.HasResponseComment {
		prComment := github.IssueComment{
			Body: &comment.CommentResponse,
		}

		if _, _, err := client.Issues.CreateComment(ctx, event.Repository.Owner, event.Repository.Name, event.PullRequest.Number, &prComment); err != nil {
			return err
		}
	}

	var apps []string
	if comment.Command.ExplicitApplication {
		apps = comment.Command.Applications
		h.Log.Debug("apps specified by comment comment %v", apps)
	} else {
		resolver := NewApplicationResolver(client, &h.ArgoClient, h.Log)
		apps, err = resolver.FindApplicationNames(ctx, event)
		h.Log.Debug("apps discovered %v", apps)
		if err != nil {
			h.Log.Debug("pull request did not change any applications managed by argocd")
			return nil
		}
	}

	commenter := vsc.NewCommenter(client, h.Log, ctx)

	if comment.Command.Name == command.Plan {
		planner := argocd.NewPlanner(&h.ArgoClient, h.Log)

		for _, app := range apps {
			var err error = nil

			h.Log.Debug("running plan for app %s against revision %s", app, event.Revision)
			plan, diff, err := planner.Plan(ctx, app, event.Revision)
			if err != nil {
				h.Log.Err(err, fmt.Sprintf("unable to plan: %s", plan))
				return err
			}

			h.Log.Debug("%s diff %t", app, diff)
			var comment string
			if diff {
				comment = fmt.Sprintf("argocd plan for `%s`\n\n", app) + "```diff\n" + plan + "\n```"
			} else {
				comment = "no diff detected, current state is up to date with this revision."
				h.Log.Info(plan)
			}

			err = commenter.Plan(&event, app, command.Plan.String(), comment)
			if err != nil {
				h.Log.Err(err, fmt.Sprintf("error while planning %s", app))
			}
		}
		return nil
	}

	if comment.Command.Name == command.Apply {
		// TODO: allow for multiple apps to be applied, I want to be careful about this
		if len(apps) != 1 {
			h.Log.Info("requested apply with more than 1 app, only one app allowed when applying")
			return nil
		}
		go func() {
			for _, app := range apps {
				applyContext, cancel := context.WithTimeout(ctx, 60*time.Second)
				defer cancel()
				apply := NewApplyRunner(client, h.Config, h.Log, &h.ArgoClient)
				response, err := apply.Run(applyContext, app, event)
				if err != nil {
					h.Log.Err(err, "unable to apply")
					return
				}
				comment := fmt.Sprintf("apply result for `%s`\n\n", app) + "```\n" + response.Message + "\n```"
				commenter.Comment(&event, &comment)
			}
		}()
		return nil
	}

	return errors.Errorf("unsupported argo command")
}

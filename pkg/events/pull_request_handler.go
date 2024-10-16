package events

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/corymurphy/argobot/pkg/argocd"
	"github.com/corymurphy/argobot/pkg/command"
	"github.com/corymurphy/argobot/pkg/env"
	vsc "github.com/corymurphy/argobot/pkg/github"
	"github.com/corymurphy/argobot/pkg/logging"
	"github.com/corymurphy/argobot/pkg/utils"
	"github.com/google/go-github/v53/github"
	"github.com/palantir/go-githubapp/githubapp"
)

var validPullActions = []string{"opened", "reopened", "ready_for_review"}

type PullRequestHandler struct {
	githubapp.ClientCreator
	Config     *env.Config
	Log        logging.SimpleLogging
	ArgoClient argocd.ApplicationsClient
}

func (h *PullRequestHandler) Handles() []string {
	return []string{"pull_request"}
}

func (h *PullRequestHandler) Handle(ctx context.Context, eventType string, deliveryID string, payload []byte) error {

	var pull github.PullRequestEvent
	if err := json.Unmarshal(payload, &pull); err != nil {
		h.Log.Err(err, "invalid github event payload")
		return fmt.Errorf("invalid github event payload")
	}

	if !utils.StringInSlice(*pull.Action, validPullActions) {
		h.Log.Debug("ignoring pull request action %s", *pull.Action)
		return nil
	}

	event := vsc.InitializeFromPullRequest(pull)
	installationID := githubapp.GetInstallationIDFromEvent(&pull)
	client, err := h.NewInstallationClient(installationID)
	if err != nil {
		return err
	}

	resolver := NewApplicationResolver(client, &h.ArgoClient, h.Log)
	apps, err := resolver.FindApplicationNames(ctx, event)
	if err != nil {
		h.Log.Debug("pull request did not change any applications managed by argocd")
		return nil
	}

	commenter := vsc.NewCommenter(client, h.Log, ctx)
	planner := argocd.NewPlanner(&h.ArgoClient, h.Log)

	for _, app := range apps {

		h.Log.Debug("running plan for app %s against revision %s", app, event.Revision)
		plan, diff, err := planner.Plan(ctx, app, event.Revision)
		if err != nil {
			h.Log.Err(err, fmt.Sprintf("unable to plan: %s", plan))
			return err
		}

		var comment string
		if diff {
			comment = fmt.Sprintf("argocd plan for `%s`\n\n", app) + "```diff\n" + plan + "\n```"
		} else {
			comment = "no diff detected, current state is up to date with this revision."
		}

		err = commenter.Plan(&event, app, command.Plan.String(), comment)
		if err != nil {
			h.Log.Err(err, fmt.Sprintf("error while planning %s", app))
		}
	}

	h.Log.Debug("ignoring unsupported event %s %s", eventType, *pull.Action)
	return nil
}

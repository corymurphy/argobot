package github

import (
	"context"

	"github.com/corymurphy/argobot/pkg/logging"
	"github.com/google/go-github/v53/github"
)

type VscClient struct {
	client *github.Client
	logger logging.SimpleLogging
}

func (v *VscClient) SetStatusCheck(ctx context.Context, event Event, state CommitState, context string, description string) error {

	url := ""
	status := &github.RepoStatus{
		State:       github.String(state.String()),
		Description: github.String(description),
		Context:     github.String(context),
		TargetURL:   &url,
	}
	_, resp, err := v.client.Repositories.CreateStatus(ctx, event.Repository.Owner, event.Repository.Name, event.Revision, status)
	if resp != nil {
		v.logger.Debug("POST /repos/%v/%v/statuses/%s returned: %v", event.Repository.Owner, event.Repository.Name, event.Revision, resp.StatusCode)
	}
	return err
}

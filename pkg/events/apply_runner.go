package events

import (
	"context"
	"fmt"

	"github.com/corymurphy/argobot/pkg/argocd"
	"github.com/corymurphy/argobot/pkg/env"
	vsc "github.com/corymurphy/argobot/pkg/github"
	"github.com/corymurphy/argobot/pkg/logging"
	"github.com/corymurphy/argobot/pkg/models"
	"github.com/google/go-github/v53/github"
)

type ApplyRunner struct {
	vcsClient   *github.Client
	Config      *env.Config
	Log         logging.SimpleLogging
	ApplyClient argocd.ApplyClient
}

func NewApplyRunner(vcsClient *github.Client, config *env.Config, log logging.SimpleLogging, applyClient argocd.ApplyClient) *ApplyRunner {
	return &ApplyRunner{
		vcsClient:   vcsClient,
		Config:      config,
		Log:         log,
		ApplyClient: applyClient,
	}
}

// TODO: validate that the PR is in an approved/mergeable state
func (a *ApplyRunner) Run(ctx context.Context, cmd *CommentCommand, comment models.PullRequestComment) (models.CommentResponse, error) {
	var resp models.CommentResponse

	status, err := vsc.NewPullRequestStatusFetcher(ctx, a.Log, a.vcsClient).Fetch(comment.Pull)
	if err != nil {
		return resp, fmt.Errorf("unable to get status of pr %w", err)
	}
	if !status.ApprovalStatus.IsApproved || !status.Mergeable {
		a.Log.Debug(
			"pull request was not in a mergeable state. approved %t, mergeable %t",
			status.ApprovalStatus.IsApproved,
			status.Mergeable)
		return models.NewCommentResponse("pull request must be approved and in a mergeable state", comment), nil
	}

	apply, err := a.ApplyClient.Apply(cmd.Application, comment.Sha)
	if err != nil {
		return resp, fmt.Errorf("argoclient failed while applying %w", err)
	}
	a.Log.Debug(apply)
	return models.NewCommentResponse(apply, comment), nil
}
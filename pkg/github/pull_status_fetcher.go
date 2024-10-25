package github

import (
	"context"

	"github.com/corymurphy/argobot/pkg/logging"
	"github.com/corymurphy/argobot/pkg/models"
	"github.com/google/go-github/v53/github"
	"github.com/pkg/errors"
)

type VscPullRequestStatusFetcher interface {
	Fetch(ctx context.Context, event Event) (models.PullRequestStatus, error)
}

type GithubPullRequestStatusFetcher struct {
	client *github.Client
	logger logging.SimpleLogging
}

func NewPullRequestStatusFetcher(logger logging.SimpleLogging, client *github.Client) VscPullRequestStatusFetcher {
	return &GithubPullRequestStatusFetcher{
		client: client,
		logger: logger,
	}
}

func (g *GithubPullRequestStatusFetcher) Fetch(ctx context.Context, event Event) (status models.PullRequestStatus, err error) {

	g.logger.Debug("getting approval status of pull request %d", event.PullRequest.Number)
	approvalStatus, err := g.PullIsApproved(ctx, event)
	if err != nil {
		return status, errors.Wrapf(err, "fetching pull approval status for repo: %s/%s, and pull number: %d", event.Repository.Owner, event.Repository.Name, event.PullRequest.Number)
	}

	mergeable, err := g.PullIsMergeable(ctx, event, "")
	if err != nil {
		return status, errors.Wrapf(err, "fetching mergeability status for repo: %s/%s, and pull number: %d", event.Repository.Owner, event.Repository.Name, event.PullRequest.Number)
	}

	return models.PullRequestStatus{
		ApprovalStatus: approvalStatus,
		Mergeable:      mergeable,
	}, err
}

func (g *GithubPullRequestStatusFetcher) PullIsApproved(ctx context.Context, event Event) (approvalStatus models.ApprovalStatus, err error) {
	g.logger.Debug("checking if pull request %d is approved", event.PullRequest.Number)
	nextPage := 0
	for {
		opts := github.ListOptions{
			PerPage: 300,
		}
		if nextPage != 0 {
			opts.Page = nextPage
		}
		pageReviews, resp, err := g.client.PullRequests.ListReviews(ctx, event.Repository.Owner, event.Repository.Name, event.PullRequest.Number, &opts)
		if resp != nil {
			g.logger.Debug("GET /repos/%v/%v/pulls/%d/reviews returned: %v", event.Repository.Owner, event.Repository.Name, event.PullRequest.Number, resp.StatusCode)
		}
		if err != nil {
			return approvalStatus, errors.Wrap(err, "getting reviews")
		}
		for _, review := range pageReviews {
			if review != nil && review.GetState() == "APPROVED" {
				return models.ApprovalStatus{
					IsApproved: true,
					ApprovedBy: *review.User.Login,
					Date:       review.SubmittedAt.Time,
				}, nil
			}
		}
		if resp.NextPage == 0 {
			break
		}
		nextPage = resp.NextPage
	}
	return approvalStatus, nil
}

func (g *GithubPullRequestStatusFetcher) PullIsMergeable(ctx context.Context, event Event, vcsstatusname string) (bool, error) {
	g.logger.Debug("Checking if GitHub pull request %d is mergeable", event.PullRequest.Number)
	githubPR, _, err := g.client.PullRequests.Get(ctx, event.Repository.Owner, event.Repository.Name, event.PullRequest.Number)
	if err != nil {
		return false, errors.Wrap(err, "getting pull request")
	}

	vscClient := NewClient(g.client, g.logger)

	state := githubPR.GetMergeableState()
	if state != "clean" && state != "unstable" && state != "has_hooks" {
		approved, err := vscClient.GetPullReviewDecision(event)
		if err != nil {
			return false, errors.Wrap(err, "getting pull request reviewDecision")
		}
		return approved, nil
	}
	return true, nil
}

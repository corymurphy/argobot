package github

import (
	"context"

	"github.com/corymurphy/argobot/pkg/logging"
	"github.com/corymurphy/argobot/pkg/models"
	"github.com/google/go-github/v53/github"
	"github.com/pkg/errors"
)

type VscPullRequestStatusFetcher interface {
	Fetch(event Event) (models.PullRequestStatus, error)
}

type GithubPullRequestStatusFetcher struct {
	client *github.Client
	ctx    context.Context
	logger logging.SimpleLogging
}

func NewPullRequestStatusFetcher(ctx context.Context, logger logging.SimpleLogging, client *github.Client) VscPullRequestStatusFetcher {
	return &GithubPullRequestStatusFetcher{
		ctx:    ctx,
		client: client,
		logger: logger,
	}
}

func (g *GithubPullRequestStatusFetcher) Fetch(event Event) (status models.PullRequestStatus, err error) {

	g.logger.Debug("getting approval status of pull request %d", event.PullRequest.Number)
	approvalStatus, err := g.PullIsApproved(event)
	if err != nil {
		return status, errors.Wrapf(err, "fetching pull approval status for repo: %s/%s, and pull number: %d", event.Repository.Owner, event.Repository.Name, event.PullRequest.Number)
	}

	mergeable, err := g.PullIsMergeable(event, "")
	if err != nil {
		return status, errors.Wrapf(err, "fetching mergeability status for repo: %s/%s, and pull number: %d", event.Repository.Owner, event.Repository.Name, event.PullRequest.Number)
	}

	return models.PullRequestStatus{
		ApprovalStatus: approvalStatus,
		Mergeable:      mergeable,
	}, err
}

func (g *GithubPullRequestStatusFetcher) PullIsApproved(event Event) (approvalStatus models.ApprovalStatus, err error) {
	g.logger.Debug("Checking if GitHub pull request %d is approved", event.PullRequest.Number)
	nextPage := 0
	for {
		opts := github.ListOptions{
			PerPage: 300,
		}
		if nextPage != 0 {
			opts.Page = nextPage
		}
		pageReviews, resp, err := g.client.PullRequests.ListReviews(g.ctx, event.Repository.Owner, event.Repository.Name, event.PullRequest.Number, &opts)
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

// TODO: complete the bypass functionality from atlantis
// TODO: check if approved by codeowner
func (g *GithubPullRequestStatusFetcher) PullIsMergeable(event Event, vcsstatusname string) (bool, error) {
	g.logger.Debug("Checking if GitHub pull request %d is mergeable", event.PullRequest.Number)
	githubPR, _, err := g.client.PullRequests.Get(g.ctx, event.Repository.Owner, event.Repository.Name, event.PullRequest.Number)
	if err != nil {
		return false, errors.Wrap(err, "getting pull request")
	}

	state := githubPR.GetMergeableState()
	if state != "clean" && state != "unstable" && state != "has_hooks" {
		//mergeable bypass apply code hidden by feature flag
		// if g.config.AllowMergeableBypassApply {
		// 	logger.Debug("AllowMergeableBypassApply feature flag is enabled - attempting to bypass apply from mergeable requirements")
		// 	if state == "blocked" {
		// 		//check status excluding atlantis apply
		// 		status, err := g.GetCombinedStatusMinusApply(logger, repo, githubPR, vcsstatusname)
		// 		if err != nil {
		// 			return false, errors.Wrap(err, "getting pull request status")
		// 		}

		// 		//check to see if pr is approved using reviewDecision
		// 		approved, err := g.GetPullReviewDecision(pull)
		// 		if err != nil {
		// 			return false, errors.Wrap(err, "getting pull request reviewDecision")
		// 		}

		// 		//if all other status checks EXCEPT atlantis/apply are successful, and the PR is approved based on reviewDecision, let it proceed
		// 		if status && approved {
		// 			return true, nil
		// 		}
		// 	}
		// }

		return false, nil
	}
	return true, nil
}

// func (g *GithubPullRequestStatusFetcher) GetPullReviewDecision(pull models.PullRequest) (approvalStatus bool, err error) {
// 	var query struct {
// 		Repository struct {
// 			PullRequest struct {
// 				ReviewDecision string
// 			} `graphql:"pullRequest(number: $number)"`
// 		} `graphql:"repository(owner: $owner, name: $name)"`
// 	}

// 	variables := map[string]interface{}{
// 		"owner":  githubv4.String(pull.Owner),
// 		"name":   githubv4.String(pull.Name),
// 		"number": githubv4.Int(pull.Number),
// 	}

// 	// g.client.BaseURL.Query()
// 	err = g.v4Client.Query(g.ctx, &query, variables)
// 	if err != nil {
// 		return approvalStatus, errors.Wrap(err, "getting reviewDecision")
// 	}

// 	if query.Repository.PullRequest.ReviewDecision == "APPROVED" || len(query.Repository.PullRequest.ReviewDecision) == 0 {
// 		return true, nil
// 	}

// 	return false, nil
// }

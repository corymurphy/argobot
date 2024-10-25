package github

// {"data":{"repository":{"pullRequest":{"reviewDecision":"APPROVED"}}}}

type PullRequestReviewDecision struct {
	Data struct {
		Repository struct {
			PullRequest struct {
				ReviewDecision string `json:"reviewDecision"`
			} `json:"pullRequest"`
		} `json:"repository"`
	} `json:"data"`
}

func (p *PullRequestReviewDecision) Approved() bool {
	return p.Data.Repository.PullRequest.ReviewDecision == "APPROVED"
}

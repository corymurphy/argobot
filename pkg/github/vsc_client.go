package github

import (
	"context"
	"encoding/json"
	"fmt"

	"io"

	"github.com/corymurphy/argobot/pkg/logging"
	"github.com/google/go-github/v53/github"
	"github.com/pkg/errors"
)

type VscClient struct {
	client *github.Client
	logger logging.SimpleLogging
}

func NewClient(client *github.Client, logger logging.SimpleLogging) *VscClient {
	return &VscClient{
		client: client,
		logger: logger,
	}
}

type reviewDecisionRequest struct {
	Query     string `json:"query"`
	Variables string `json:"variables"`
}

func (v *VscClient) GetPullReviewDecision(event Event) (approvalStatus bool, err error) {

	query := `
		query($owner: String!, $name: String!, $number: Int!) {
			repository(owner: $owner, name: $name) {
				pullRequest(number: $number) {
					reviewDecision
				}
			}
		}
		
	`
	variables := fmt.Sprintf(`
		{
			"owner": "%s",
			"name": "%s",
			"number": %d
		}
	`, event.Repository.Owner, event.Repository.Name, event.PullRequest.Number)

	body := reviewDecisionRequest{
		Query:     query,
		Variables: variables,
	}

	request, err := v.client.NewRequest("POST", "/graphql", body)

	if err != nil {
		return false, errors.Wrap(err, "unable to create http request")
	}

	response, err := v.client.Client().Do(request)
	if err != nil {
		return false, errors.Wrap(err, "failed while calling the graphql api")
	}
	defer response.Body.Close()
	dataSerialized, err := io.ReadAll(io.Reader(response.Body))
	if err != nil {
		return false, errors.Wrap(err, "unable to read response")
	}
	v.logger.Info(string(dataSerialized))

	var decision PullRequestReviewDecision
	err = json.Unmarshal(dataSerialized, &decision)
	if err != nil {
		return false, errors.Wrap(err, "unable to deserialize response")
	}
	return decision.Approved(), nil
}

func (v *VscClient) SetStatusCheck(ctx context.Context, event Event, state CommitState, context string, description string, url string) error {

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

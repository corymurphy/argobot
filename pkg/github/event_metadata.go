package github

import (
	"encoding/json"
	"fmt"

	"github.com/google/go-github/v53/github"
	"github.com/pkg/errors"
)

type Action int

const (
	Comment Action = iota
	Opened  Action = iota
)

func (e Action) String() string {
	switch e {
	case Comment:
		return "comment"
	case Opened:
		return "opened"
	}
	return ""
}

type Repository struct {
	Name  string
	Owner string
}

type Actor struct {
	Name string
}

type PullRequest struct {
	Number int
}

type GithubEvent struct {
	Action *string `json:"action,omitempty"`
}

type InstallationProvider interface {
	GetInstallation() *github.Installation
}

type Event struct {
	Revision             string
	Message              string
	IsPullRequest        bool
	Action               Action
	Actor                Actor
	Repository           Repository
	PullRequest          PullRequest
	InstallationProvider InstallationProvider
}

func (e *Event) GetInstallation() *github.Installation {
	return e.InstallationProvider.GetInstallation()
}

func (e *Event) HasMessage() bool {
	return e.Message != ""
}

func NewEvent(eventType string, payload []byte) (Event, error) {
	var event Event
	var githubEvent GithubEvent
	if err := json.Unmarshal(payload, &githubEvent); err != nil {
		return event, errors.Wrap(err, "failed to parse event payload")
	}

	if eventType == "issue_comment" && *githubEvent.Action == "created" {
		var comment github.IssueCommentEvent
		if err := json.Unmarshal(payload, &comment); err != nil {
			return event, errors.Wrap(err, "failed to parse issue comment event payload")
		}

		return Event{
			Actor:         Actor{Name: comment.GetComment().GetUser().GetLogin()},
			Action:        Comment,
			IsPullRequest: comment.GetIssue().IsPullRequest(),
			Repository: Repository{
				Name:  comment.GetRepo().GetName(),
				Owner: comment.GetRepo().GetOwner().GetLogin(),
			},
			PullRequest: PullRequest{
				Number: comment.GetIssue().GetNumber(),
			},
			Message:              *comment.GetComment().Body,
			InstallationProvider: &comment,
		}, nil
	}

	if eventType == "pull_request" && (*githubEvent.Action == "opened" || *githubEvent.Action == "reopened" || *githubEvent.Action == "ready_for_review") {
		var pr github.PullRequestEvent
		if err := json.Unmarshal(payload, &pr); err != nil {
			return event, errors.Wrap(err, "failed to parse issue comment event payload")
		}

		return Event{
			Actor:         Actor{Name: pr.GetPullRequest().GetUser().GetLogin()},
			Action:        Opened,
			IsPullRequest: true,
			Repository: Repository{
				Name:  pr.GetRepo().GetName(),
				Owner: pr.GetRepo().GetOwner().GetLogin(),
			},
			PullRequest: PullRequest{
				Number: pr.GetPullRequest().GetNumber(),
			},
			Message:              "",
			InstallationProvider: &pr,
		}, nil
	}

	return event, fmt.Errorf("unsupported event %s %s", eventType, *githubEvent.Action)
}

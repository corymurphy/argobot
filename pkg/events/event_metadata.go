package events

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

type InstallationProvider interface {
	GetInstallation() *github.Installation
}

type EventMetadata struct {
	Revision             string
	Message              string
	IsPullRequest        bool
	Action               Action
	Actor                Actor
	Repository           Repository
	PullRequest          PullRequest
	InstallationProvider InstallationProvider
}

func (e *EventMetadata) GetInstallation() *github.Installation {
	return e.InstallationProvider.GetInstallation()
}

func (e *EventMetadata) HasMessage() bool {
	return e.Message != ""
}

func NewEventMetadata(eventType string, payload []byte) (EventMetadata, error) {
	var metadata EventMetadata
	var event Event
	if err := json.Unmarshal(payload, &event); err != nil {
		return metadata, errors.Wrap(err, "failed to parse event payload")
	}

	if eventType == "issue_comment" && *event.Action == "created" {
		var event github.IssueCommentEvent
		if err := json.Unmarshal(payload, &event); err != nil {
			return metadata, errors.Wrap(err, "failed to parse issue comment event payload")
		}

		return EventMetadata{
			Actor:         Actor{Name: event.GetComment().GetUser().GetLogin()},
			Action:        Comment,
			IsPullRequest: event.GetIssue().IsPullRequest(),
			Repository: Repository{
				Name:  event.GetRepo().GetName(),
				Owner: event.GetRepo().GetOwner().GetLogin(),
			},
			PullRequest: PullRequest{
				Number: event.GetIssue().GetNumber(),
			},
			Message:              *event.GetComment().Body,
			InstallationProvider: &event,
		}, nil
	}

	if eventType == "pull_request" && (*event.Action == "opened" || *event.Action == "reopened" || *event.Action == "ready_for_review") {
		var event github.PullRequestEvent
		if err := json.Unmarshal(payload, &event); err != nil {
			return metadata, errors.Wrap(err, "failed to parse issue comment event payload")
		}

		return EventMetadata{
			Actor:         Actor{Name: event.GetPullRequest().GetUser().GetLogin()},
			Action:        Opened,
			IsPullRequest: true,
			Repository: Repository{
				Name:  event.GetRepo().GetName(),
				Owner: event.GetRepo().GetOwner().GetLogin(),
			},
			PullRequest: PullRequest{
				Number: event.GetPullRequest().GetNumber(),
			},
			Message:              "",
			InstallationProvider: &event,
		}, nil
	}

	return metadata, fmt.Errorf("unsupported event %s %s", eventType, *event.Action)
}

type GithubEvent struct {
	// Action is the action that was performed on the comment.
	// Possible values are: "created", "edited", "deleted".
	Action *string `json:"action,omitempty"`
}

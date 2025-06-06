package github

import (
	"fmt"

	"github.com/google/go-github/v72/github"
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

func (r *Repository) HtmlUrl() string {
	return fmt.Sprintf("https://github.com/%s/%s", r.Owner, r.Name)
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
	GithubClient         *github.Client
}

func (e *Event) GetInstallation() *github.Installation {
	return e.InstallationProvider.GetInstallation()
}

func (e *Event) HasMessage() bool {
	return e.Message != ""
}

func InitializeFromIssueComment(source github.IssueCommentEvent, revision string) (Event, *int64) {
	return Event{
		Actor:         Actor{Name: source.GetComment().GetUser().GetLogin()},
		Action:        Comment,
		IsPullRequest: source.GetIssue().IsPullRequest(),
		Revision:      revision,
		Repository: Repository{
			Name:  source.GetRepo().GetName(),
			Owner: source.GetRepo().GetOwner().GetLogin(),
		},
		PullRequest: PullRequest{
			Number: source.GetIssue().GetNumber(),
		},
		Message:              *source.GetComment().Body,
		InstallationProvider: &source,
	}, source.Comment.ID
}

func InitializeFromPullRequest(source github.PullRequestEvent) Event {
	return Event{
		Actor:         Actor{Name: source.GetPullRequest().GetUser().GetLogin()},
		Action:        Opened,
		IsPullRequest: true,
		Revision:      *source.PullRequest.Head.SHA,
		Repository: Repository{
			Name:  source.GetRepo().GetName(),
			Owner: source.GetRepo().GetOwner().GetLogin(),
		},
		PullRequest: PullRequest{
			Number: source.GetPullRequest().GetNumber(),
		},
		Message: "",
	}
}

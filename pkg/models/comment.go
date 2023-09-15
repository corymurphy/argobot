package models

type PullRequestComment struct {
	Body string
}

func NewPullRequestComment(body string) PullRequestComment {
	return PullRequestComment{
		Body: body,
	}
}

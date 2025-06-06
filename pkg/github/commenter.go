package github

import (
	"context"
	"fmt"

	"github.com/corymurphy/argobot/pkg/logging"
	"github.com/corymurphy/argobot/pkg/utils"
	"github.com/google/go-github/v72/github"
)

const maxCommentLength = 32768

type VscCommenter interface {
	Plan(event *Event, app string, command string, comment string)
	Comment(event *Event, comment *string)
}

type Commenter struct {
	client *github.Client
	log    logging.SimpleLogging
	ctx    context.Context
}

func NewCommenter(client *github.Client, log logging.SimpleLogging, ctx context.Context) *Commenter {
	return &Commenter{
		client: client,
		log:    log,
		ctx:    ctx,
	}
}

func (c *Commenter) Plan(event *Event, app string, command string, comment string) error {
	c.log.Debug("Creating comment on GitHub pull request %d", event.PullRequest.Number)
	var sepStart string

	sepEnd := "\n```\n</details>" +
		"\n<br>\n\n**Warning**: Output length greater than max comment size. Continued in next comment."

	if command != "" {
		sepStart = fmt.Sprintf("Continued %s output from previous comment.\n<details><summary>Show Output</summary>\n\n", command) +
			"```diff\n"
	} else {
		sepStart = "Continued from previous comment.\n<details><summary>Show Output</summary>\n\n" +
			"```diff\n"
	}

	truncationHeader := "\n```\n</details>" +
		"\n<br>\n\n**Warning**: Command output is larger than the maximum number of comments per command. Output truncated.\n\n[..]\n"

	comments := utils.SplitComment(comment, maxCommentLength, sepEnd, sepStart, 100, truncationHeader)
	for i := range comments {
		_, resp, err := c.client.Issues.CreateComment(c.ctx, event.Repository.Owner, event.Repository.Name, event.PullRequest.Number, &github.IssueComment{Body: &comments[i]})
		if resp != nil {
			c.log.Debug("POST /repos/%v/%v/issues/%d/comments returned: %v", event.Repository.Owner, event.Repository.Name, event.PullRequest.Number, resp.StatusCode)
		}
		if err != nil {
			return err
		}
	}

	return nil
}

func (c *Commenter) Comment(event *Event, comment *string) error {
	_, resp, err := c.client.Issues.CreateComment(c.ctx, event.Repository.Owner, event.Repository.Name, event.PullRequest.Number, &github.IssueComment{Body: comment})
	if resp != nil {
		c.log.Debug("POST /repos/%v/%v/issues/%d/comments returned: %v", event.Repository.Owner, event.Repository.Name, event.PullRequest.Number, resp.StatusCode)
	}
	if err != nil {
		return err
	}
	return err
}

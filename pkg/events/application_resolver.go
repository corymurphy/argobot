package events

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/corymurphy/argobot/pkg/argocd"
	"github.com/corymurphy/argobot/pkg/logging"
	"github.com/corymurphy/argobot/pkg/models"
	"github.com/google/go-github/v53/github"
)

type ApplicationResolver struct {
	GitHubClient *github.Client
	ArgoCdClient *argocd.ApplicationsClient
	Log          logging.SimpleLogging
}

func NewApplicationResolver(githubClient *github.Client, argocdClient *argocd.ApplicationsClient, log logging.SimpleLogging) *ApplicationResolver {
	return &ApplicationResolver{
		GitHubClient: githubClient,
		ArgoCdClient: argocdClient,
		Log:          log,
	}
}

func (a *ApplicationResolver) FindApplicationNames(ctx context.Context, command *CommentCommand, pull models.PullRequest) ([]string, error) {
	var changedApps []string

	modified, err := a.GetModifiedFiles(ctx, pull)
	if err != nil {
		return changedApps, fmt.Errorf("unable to list github pull request modified files: %s", err)
	}

	apps, err := a.ArgoCdClient.List()
	if err != nil {
		return changedApps, fmt.Errorf("unable to list argo apps: %s", err)
	}

	for _, app := range apps.Items {

		path := app.Spec.Source.Path
		name := app.Name
		// name := app.Metadata.Name
		// app.Name
		a.Log.Debug(fmt.Sprintf("name: %s | path: %s", name, path))

		for _, file := range modified {
			if strings.Contains(file, path) {
				changedApps = append(changedApps, name)
			}
		}
	}

	return changedApps, nil
}

// copied from atlantis
func (a *ApplicationResolver) GetModifiedFiles(ctx context.Context, pull models.PullRequest) ([]string, error) {
	a.Log.Debug("Getting modified files for GitHub pull request %d", pull.Number)
	var files []string
	nextPage := 0

listloop:
	for {
		opts := github.ListOptions{
			PerPage: 300,
		}
		if nextPage != 0 {
			opts.Page = nextPage
		}
		// GitHub has started to return 404's sometimes. They've got some
		// eventual consistency issues going on so we're just going to attempt
		// up to 5 times for each page with exponential backoff.
		maxAttempts := 5
		attemptDelay := 0 * time.Second
		for i := 0; i < maxAttempts; i++ {
			// First don't sleep, then sleep 1, 3, 7, etc.
			time.Sleep(attemptDelay)
			attemptDelay = 2*attemptDelay + 1*time.Second

			pageFiles, resp, err := a.GitHubClient.PullRequests.ListFiles(ctx, pull.Owner, pull.Name, pull.Number, &opts)
			if resp != nil {
				a.Log.Debug("[attempt %d] GET /repos/%v/%v/pulls/%d/files returned: %v", i+1, pull.Owner, pull.Name, pull.Number, resp.StatusCode)
			}
			if err != nil {
				ghErr, ok := err.(*github.ErrorResponse)
				if ok && ghErr.Response.StatusCode == 404 {
					// (hopefully) transient 404, retry after backoff
					continue
				}
				// something else, give up
				return files, err
			}
			for _, f := range pageFiles {
				files = append(files, f.GetFilename())

				// If the file was renamed, we'll want to run plan in the directory
				// it was moved from as well.
				if f.GetStatus() == "renamed" {
					files = append(files, f.GetPreviousFilename())
				}
			}
			if resp.NextPage == 0 {
				break listloop
			}
			nextPage = resp.NextPage
			break
		}
	}
	return files, nil
}

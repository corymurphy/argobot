package events

import (
	"github.com/corymurphy/argobot/pkg/argocd"
	"github.com/corymurphy/argobot/pkg/env"
	"github.com/corymurphy/argobot/pkg/logging"
	"github.com/palantir/go-githubapp/githubapp"
)

type PullRequestHandler struct {
	githubapp.ClientCreator
	Config      *env.Config
	Log         logging.SimpleLogging
	ArgoClient  argocd.ApplicationsClient
	TestingMode bool
}

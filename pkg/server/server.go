package server

import (
	"fmt"
	"net/http"
	"time"

	"github.com/corymurphy/argobot/pkg/argocd"
	"github.com/corymurphy/argobot/pkg/env"
	"github.com/corymurphy/argobot/pkg/events"
	"github.com/corymurphy/argobot/pkg/logging"
	"github.com/gregjones/httpcache"
	"github.com/palantir/go-githubapp/githubapp"
)

type Server struct {
	Config           *env.Config
	Log              logging.SimpleLogging
	PRCommentHandler http.Handler
	PlanClient       argocd.PlanClient
	ApplyClient      argocd.ApplyClient
	githubapp.ClientCreator
}

func NewServer(config *env.Config, logger logging.SimpleLogging, planClient argocd.PlanClient, applyClient argocd.ApplyClient) *Server {

	cc, err := githubapp.NewDefaultCachingClientCreator(
		config.Github,
		githubapp.WithClientUserAgent("argobot/0.1.0"),
		githubapp.WithClientTimeout(3*time.Second),
		githubapp.WithClientCaching(false, func() httpcache.Cache { return httpcache.NewMemoryCache() }),
	)

	if err != nil {
		logger.Err("unable to create github client creator", err)
		panic(err)
	}

	prCommentHandler := &events.PRCommentHandler{
		ClientCreator: cc,
		Config:        config,
		Log:           logger,
		PlanClient:    planClient,
		ApplyClient:   applyClient,
	}

	webhookHandler := githubapp.NewDefaultEventDispatcher(config.Github, prCommentHandler)

	return &Server{
		Config:           config,
		Log:              logger,
		ClientCreator:    cc,
		PRCommentHandler: webhookHandler,
		PlanClient:       planClient,
		ApplyClient:      applyClient,
	}
}

func (s *Server) Health(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	w.Write([]byte(`{
  "status": "ok"
}`))
}

func (s *Server) Start() {

	http.Handle("/health", http.HandlerFunc(s.Health))
	http.Handle(githubapp.DefaultWebhookRoute, s.PRCommentHandler)

	addr := fmt.Sprintf("%s:%d", s.Config.Server.Address, s.Config.Server.Port)
	s.Log.Info("starting server on %s", addr)

	if err := http.ListenAndServe(addr, nil); err != nil {
		panic(err)
	}
}

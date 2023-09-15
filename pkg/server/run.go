package server

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/corymurphy/argobot/pkg/env"
	"github.com/corymurphy/argobot/pkg/events"
	"github.com/gregjones/httpcache"
	"github.com/palantir/go-githubapp/githubapp"
	"github.com/rcrowley/go-metrics"
	"github.com/rs/zerolog"
)

func Start(opts Options) {
	config, err := env.ReadConfig(opts.Path)
	if err != nil {
		panic(err)
	}

	content, err := os.ReadFile(config.AppConfig.PrivateKeyFilePath)

	if err != nil {
		panic(err)
	}

	config.Github.App.PrivateKey = string(content)
	config.Github.App.WebhookSecret = os.Getenv("ARGOBOT_GH_WEBHOOK_SECRET")

	// logger := zerolog.New(os.Stdout).With().Timestamp().Logger()
	logger := zerolog.New(os.Stdout).Level(zerolog.DebugLevel).With().Timestamp().Logger()

	zerolog.DefaultContextLogger = &logger

	metricsRegistry := metrics.DefaultRegistry

	cc, err := githubapp.NewDefaultCachingClientCreator(
		config.Github,
		githubapp.WithClientUserAgent("argobot/0.1.0"),
		githubapp.WithClientTimeout(3*time.Second),
		githubapp.WithClientCaching(false, func() httpcache.Cache { return httpcache.NewMemoryCache() }),
		githubapp.WithClientMiddleware(
			githubapp.ClientMetrics(metricsRegistry),
		),
	)
	if err != nil {
		panic(err)
	}

	prCommentHandler := &events.PRCommentHandler{
		ClientCreator: cc,
		Config:        config,
	}

	webhookHandler := githubapp.NewDefaultEventDispatcher(config.Github, prCommentHandler)

	http.Handle("/health", http.HandlerFunc(healthCheck))
	http.Handle(githubapp.DefaultWebhookRoute, webhookHandler)

	addr := fmt.Sprintf("%s:%d", config.Server.Address, config.Server.Port)
	logger.Info().Msgf("Starting server on %s...", addr)
	err = http.ListenAndServe(addr, nil)
	if err != nil {
		panic(err)
	}
}

func healthCheck(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(``))
}

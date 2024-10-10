package events

import (
	"context"
	"net/http"

	"github.com/corymurphy/argobot/pkg/logging"
)

type EventHandler struct {
	Log      logging.SimpleLogging
	Handlers map[string]http.Handler
}

func (e *EventHandler) Handles() []string {
	return []string{"issue_comment", "pull_request"}
}

func (e *EventHandler) Handle(ctx context.Context, eventType string, deliveryID string, payload []byte) error {
	// return e.Handlers[eventType]

	// handler := e.Handlers[eventType]

	return nil
}

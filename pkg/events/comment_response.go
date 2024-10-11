package events

import "github.com/corymurphy/argobot/pkg/github"

type CommentResponse struct {
	Message string
	Event   github.Event
}

func NewCommentResponse(msg string, event github.Event) CommentResponse {
	return CommentResponse{
		Message: msg,
		Event:   event,
	}
}

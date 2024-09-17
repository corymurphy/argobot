package events

import "github.com/corymurphy/argobot/pkg/command"

type CommentCommand struct {
	Flags       []string
	Name        command.Name
	Application string
}

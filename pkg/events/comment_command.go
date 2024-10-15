package events

import "github.com/corymurphy/argobot/pkg/command"

func NewAutoRunCommand(apps []string, name command.Name) *CommentCommand {
	return &CommentCommand{
		Flags:               []string{},
		Name:                name,
		Applications:        apps,
		ExplicitApplication: false,
	}
}

type CommentCommand struct {
	Flags               []string
	Name                command.Name
	Applications        []string
	ExplicitApplication bool
}

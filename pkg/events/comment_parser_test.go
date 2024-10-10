package events

import (
	"testing"

	"github.com/corymurphy/argobot/pkg/command"
	"github.com/corymurphy/argobot/pkg/logging"
)

func Test_Comment_IsHelp(t *testing.T) {
	serialized := `
{
	"action": "created",
	"issue": {
		"pull_request": {
			"url": "https://api.github.com/example/example"
		}
	},
	"comment" : {
		"body": "argo help",
		"user": {
			"login": "githubuser"
		}
	}
}
`
	event, err := NewEventMetadata("issue_comment", []byte(serialized))
	if err != nil {
		t.Error(err)
	}
	parser := NewCommentParser(logging.NewLogger(logging.Silent))

	result := parser.Parse(event)
	if result.Command.Name != command.Help {
		t.Errorf("expected %s, got %s", command.Help, &result.Command.Name)
	}
}

func Test_Comment_IsBot(t *testing.T) {
	serialized := `
{
	"action": "created",
	"issue": {
		"pull_request": {
			"url": "https://api.github.com/example/example"
		}
	},
	"comment" : {
		"body": "heres the diff output",
		"user": {
			"login": "[bot] githubuser"
		}
	}
}
`
	event, err := NewEventMetadata("issue_comment", []byte(serialized))
	if err != nil {
		t.Error(err)
	}
	parser := NewCommentParser(logging.NewLogger(logging.Silent))

	result := parser.Parse(event)

	if !result.Ignore && !result.ImmediateResponse {
		t.Error("expected ignore and immediate response")
	}
}

func Test_PlanHasApplicationName(t *testing.T) {
	serialized := `
	{
		"action": "created",
		"issue": {
			"pull_request": {
				"url": "https://api.github.com/example/example"
			}
		},
		"comment" : {
			"body": "argo plan --application myapp",
			"user": {
				"login": "githubuser"
			}
		}
	}
	`

	event, err := NewEventMetadata("issue_comment", []byte(serialized))
	if err != nil {
		t.Error(err)
	}
	parser := NewCommentParser(logging.NewLogger(logging.Silent))

	result := parser.Parse(event)
	if result.Command.Name != command.Plan {
		t.Errorf("expected %s, got %s", command.Help, &result.Command.Name)
	}

	if result.Command.Application != "myapp" {
		t.Log(result.Command)
		t.Errorf("expected %s, got %s", "myapp", result.Command.Application)
	}
}

func Test_ApplyHasApplicationName(t *testing.T) {
	serialized := `
	{
		"action": "created",
		"issue": {
			"pull_request": {
				"url": "https://api.github.com/example/example"
			}
		},
		"comment" : {
			"body": "argo apply --application myapp",
			"user": {
				"login": "githubuser"
			}
		}
	}
	`

	event, err := NewEventMetadata("issue_comment", []byte(serialized))
	if err != nil {
		t.Error(err)
	}
	parser := NewCommentParser(logging.NewLogger(logging.Silent))

	result := parser.Parse(event)
	if result.Command.Name != command.Apply {
		t.Errorf("expected %s, got %s", command.Apply, &result.Command.Name)
	}

	if result.Command.Application != "myapp" {
		t.Log(result.Command)
		t.Errorf("expected %s, got %s", "myapp", result.Command.Application)
	}
}

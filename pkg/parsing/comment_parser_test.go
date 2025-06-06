package parsing

import (
	"encoding/json"
	"testing"

	"github.com/google/go-github/v72/github"
)

func Test_Comment_IsDiff(t *testing.T) {
	serialized := `
{	
	"comment" : {
		"body": "argo diff --application helloworld",
		"user": {
			"login": "githubuser"
		}
	}
}
`
	var event github.IssueCommentEvent
	json.Unmarshal([]byte(serialized), &event)
	parser, _ := NewPRCommentParser(event)

	if Diff != parser.Command {
		t.Errorf("expected %s, got %s", Diff, parser.Command)
	}
}

func Test_Comment_IsHelp(t *testing.T) {
	serialized := `
{	
	"comment" : {
		"body": "argo help",
		"user": {
			"login": "githubuser"
		}
	}
}
`
	var event github.IssueCommentEvent
	json.Unmarshal([]byte(serialized), &event)
	parser, _ := NewPRCommentParser(event)

	if Help != parser.Command {
		t.Errorf("expected %s, got %s", Diff, parser.Command)
	}
}

func Test_Comment_NotArgoCommand(t *testing.T) {
	serialized := `
{	
	"comment" : {
		"body": "test comment",
		"user": {
			"login": "githubuser"
		}
	}
}
`
	var event github.IssueCommentEvent
	json.Unmarshal([]byte(serialized), &event)
	parser, _ := NewPRCommentParser(event)

	if parser.IsArgoCommand {
		t.Error("argo command was not expected")
	}
}

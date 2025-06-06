package parsing

import (
	"errors"
	"strings"
	"unicode"

	"github.com/google/go-github/v72/github"
)

const (
	BotCommand string  = "argo"
	Diff       Command = "diff"
	Sync       Command = "sync"
	Help       Command = "help"
	Unknown    Command = "unknown"
	Error      Command = "error"
	// ArgApp        string  = "--application"
	CommandMaxLen int = 256
	AppNameMaxLen int = 128
	// ArgDir     string  = "--directory"
)

type Command string

type PRCommentParser struct {
	Event         github.IssueCommentEvent
	IsBot         bool
	IsArgoCommand bool
	Command       Command
	// Application   string
	// Directory     string
}

func IsLetter(s string) bool {
	for _, r := range s {
		if !unicode.IsLetter(r) {
			return false
		}
	}
	return true
}

func NewPRCommentParser(event github.IssueCommentEvent) (PRCommentParser, error) {

	author := event.GetComment().GetUser().GetLogin()
	isBot := strings.HasSuffix(author, "[bot]")

	if isBot {
		return PRCommentParser{
			Event: event,
			IsBot: isBot,
		}, nil
	}

	body := event.GetComment().Body

	if len(*body) > CommandMaxLen {
		return PRCommentParser{}, errors.New("pull request comment body exceeded max length")
	}

	isArgo := strings.HasPrefix(*body, BotCommand)

	if !isArgo {
		return PRCommentParser{
			Event:         event,
			IsBot:         isBot,
			IsArgoCommand: isArgo,
		}, nil
	}

	if *body == BotCommand {
		return PRCommentParser{
			Event:         event,
			IsBot:         isBot,
			IsArgoCommand: isArgo,
			Command:       Help,
		}, nil
	}

	args := strings.Split(*body, " ")

	if len(args) < 2 {
		return PRCommentParser{
			Event:         event,
			IsBot:         isBot,
			IsArgoCommand: isArgo,
			Command:       Help,
		}, nil
	}

	command := Unknown

	switch args[1] {
	case string(Diff):
		command = Diff
	case string(Help):
		command = Help
	}

	return PRCommentParser{
		Event:         event,
		IsBot:         isBot,
		IsArgoCommand: isArgo,
		Command:       command,
	}, nil
}

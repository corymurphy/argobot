package parsing

import (
	"strings"

	"github.com/google/go-github/v53/github"
)

const (
	BotCommand string  = "argo"
	Diff       Command = "diff"
	Sync       Command = "sync"
	Help       Command = "help"
	Unknown    Command = "unknown"
	Error      Command = "error"
	ArgApp     string  = "--application"
	ArgDir     string  = "--directory"
)

type Command string

type PRCommentParser struct {
	Event         github.IssueCommentEvent
	IsBot         bool
	IsArgoCommand bool
	Command       Command
	Application   string
	Directory     string
}

func NewPRCommentParser(event github.IssueCommentEvent) PRCommentParser {

	author := event.GetComment().GetUser().GetLogin()
	isBot := strings.HasSuffix(author, "[bot]")

	if isBot {
		return PRCommentParser{
			Event: event,
			IsBot: isBot,
		}
	}

	body := event.GetComment().Body
	isArgo := strings.HasPrefix(*body, BotCommand)

	if !isArgo {
		return PRCommentParser{
			Event:         event,
			IsBot:         isBot,
			IsArgoCommand: isArgo,
		}
	}

	if *body == BotCommand {
		return PRCommentParser{
			Event:         event,
			IsBot:         isBot,
			IsArgoCommand: isArgo,
			Command:       Help,
		}
	}

	args := strings.Split(*body, " ")

	if len(args) < 2 {
		return PRCommentParser{
			Event:         event,
			IsBot:         isBot,
			IsArgoCommand: isArgo,
			Command:       Help,
		}
	}

	command := Unknown

	switch args[1] {
	case string(Diff):
		command = Diff
	case string(Help):
		command = Help
	}

	if command == Unknown {
		return PRCommentParser{
			Event:         event,
			IsBot:         isBot,
			IsArgoCommand: isArgo,
			Command:       Unknown,
		}
	}

	app := ""
	dir := ""

	if command == Diff {
		for i := 2; i < len(args); i++ {

			if args[i] == ArgApp && i < len(args)+1 {

				app = args[i+1]
				i++
				continue
			}

			if args[i] == ArgDir && i < len(args)+1 {
				dir = args[i+1]
				i++
				continue
			}

		}

		if app == "" || dir == "" {
			command = Error
		}

		return PRCommentParser{
			Event:         event,
			IsBot:         isBot,
			IsArgoCommand: isArgo,
			Command:       command,
			Application:   app,
			Directory:     dir,
		}
	}

	return PRCommentParser{
		Event:         event,
		IsBot:         isBot,
		IsArgoCommand: isArgo,
		Command:       Unknown,
	}
}

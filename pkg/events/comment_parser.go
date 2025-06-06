package events

import (
	"fmt"
	"io"
	"strings"
	"unicode"

	"github.com/corymurphy/argobot/pkg/command"
	vsc "github.com/corymurphy/argobot/pkg/github"
	"github.com/corymurphy/argobot/pkg/logging"
	"github.com/google/go-github/v72/github"
	"github.com/google/shlex"
	"github.com/spf13/pflag"
)

const (
	CommandMaxLen        int    = 256
	AppNameMaxLen        int    = 128
	CommandName          string = "argo"
	applicationFlagLong  string = "application"
	applicationFlagShort string = "a"
)

type CommentParsing interface {
	Parse(event github.IssueCommentEvent) CommentParseResult
}

type CommentParser struct {
	AllowCommands []command.Name
	Log           logging.SimpleLogging
}

func NewCommentParser(logger logging.SimpleLogging) *CommentParser {
	return &CommentParser{
		Log:           logger,
		AllowCommands: []command.Name{command.Help, command.Plan, command.Apply, command.Unlock},
	}
}

func (c *CommentParser) Parse(event vsc.Event) *CommentParseResult {

	if event.Action == vsc.Opened {
		c.Log.Info("responding to pull request open event with help message")
		return &CommentParseResult{
			ImmediateResponse:  true,
			Ignore:             false,
			HasResponseComment: true,
			CommentResponse:    helpComment(),
			Command: &CommentCommand{
				Name: command.Help,
			},
		}
	}

	if event.Action != vsc.Comment {
		msg := "ignoring a non created event"
		c.Log.Info(msg)
		return &CommentParseResult{
			ImmediateResponse:  true,
			Ignore:             true,
			HasResponseComment: false,
		}
	}

	if !event.IsPullRequest {
		msg := "ignoring a non pull request event"
		c.Log.Info(msg)
		return &CommentParseResult{
			ImmediateResponse:  true,
			Ignore:             true,
			HasResponseComment: false,
		}
	}

	isBot := strings.HasSuffix(event.Actor.Name, "[bot]")
	if isBot {
		c.Log.Info("ignoring comment from a bot")
		return &CommentParseResult{
			ImmediateResponse:  true,
			Ignore:             true,
			HasResponseComment: false,
		}
	}

	body := event.Message
	if len(body) > CommandMaxLen {
		c.Log.Info("ignoring comment exceeding max length")
		return &CommentParseResult{
			ImmediateResponse:  true,
			Ignore:             true,
			HasResponseComment: false,
		}
	}

	isArgo := strings.HasPrefix(body, CommandName)
	if !isArgo {
		c.Log.Debug("ignoring, comment is not an argobot command")
		return &CommentParseResult{
			ImmediateResponse:  true,
			Ignore:             true,
			HasResponseComment: false,
		}
	}

	body = strings.ToLower(body)

	if body == CommandName || body == fmt.Sprintf("%s help", CommandName) {
		c.Log.Debug("responding, help requested")
		return &CommentParseResult{
			ImmediateResponse:  true,
			Ignore:             false,
			HasResponseComment: true,
			CommentResponse:    helpComment(),
			Command: &CommentCommand{
				Name: command.Help,
			},
		}
	}

	args, err := shlex.Split(body)
	if err != nil || len(args) < 1 {
		return &CommentParseResult{
			ImmediateResponse:  true,
			Ignore:             false,
			HasResponseComment: true,
			CommentResponse:    fmt.Sprintf("```\nError parsing command: %s\n```", err),
		}
	}

	cmd := args[1]
	if !c.isAllowedCommand(cmd) {
		c.Log.Debug("requested command %s is not allowed", cmd)
		return &CommentParseResult{
			ImmediateResponse:  true,
			Ignore:             false,
			HasResponseComment: true,
			CommentResponse:    helpComment(),
			Command: &CommentCommand{
				Name: command.Help,
			},
		}
	}

	var app string
	var flagSet *pflag.FlagSet
	var name command.Name

	switch cmd {
	case command.Plan.String():
		name = command.Plan
		flagSet = pflag.NewFlagSet(command.Plan.String(), pflag.ContinueOnError)
		flagSet.SetOutput(io.Discard)
		flagSet.StringVarP(&app, applicationFlagLong, applicationFlagShort, "", "ArgoCD application to run plan against")
	case command.Apply.String():
		name = command.Apply
		flagSet = pflag.NewFlagSet(command.Apply.String(), pflag.ContinueOnError)
		flagSet.SetOutput(io.Discard)
		flagSet.StringVarP(&app, applicationFlagLong, applicationFlagShort, "", "ArgoCD application to apply changes to")
	case command.Unlock.String():
		name = command.Unlock
		flagSet = pflag.NewFlagSet(command.Unlock.String(), pflag.ContinueOnError)
		flagSet.SetOutput(io.Discard)
		flagSet.StringVarP(&app, applicationFlagLong, applicationFlagShort, "", "ArgoCD application to unlock")
	default:
		c.Log.Debug("failed to parse command %s", cmd)
		return &CommentParseResult{
			Ignore:            true,
			ImmediateResponse: true,
		}
	}

	flagSet.Parse(args[2:])

	if app == "" {
		return &CommentParseResult{
			Command: &CommentCommand{
				Name:                name,
				Applications:        []string{},
				ExplicitApplication: false,
			},
		}
	}

	if !c.isValidApplicationName(app) {
		c.Log.Debug("invalid application name %s", app)
		return &CommentParseResult{
			Ignore:            true,
			ImmediateResponse: true,
		}
	}

	return &CommentParseResult{
		Command: &CommentCommand{
			Name:                name,
			Applications:        []string{app},
			ExplicitApplication: true,
		},
	}
}

func (e *CommentParser) isValidApplicationName(name string) bool {
	for _, s := range name {
		if !unicode.IsLetter(s) && s != '-' && !unicode.IsDigit(s) {
			return false
		}
	}
	return true
}

func (e *CommentParser) isAllowedCommand(cmd string) bool {
	for _, allowed := range e.AllowCommands {
		if allowed.String() == cmd {
			return true
		}
	}
	return false
}

func helpComment() string {
	return "```shell\n" + `
Usage: argo COMMAND [ARGS]...

  Allows you to interact with ArgoCD from a Pull Request.

  If the application name is not provided, it will be discovered
  based on the location of the modified files.

  Note: Please use care when running the apply command.
	It does not lock the state between plans like atlantis right now.

Examples:
  # show the argobot help message
  argo help

  # show the plan for all modified applications
  argo plan

  # show the plan for a specific application
  argo plan --application myapp

  # apply changes for a specific application
  argo apply --application myapp

  # unlock all locked applications modified by this pull request
  argo unlock

Commands:
  help 		Shows this message
  plan 		[ --application myapp ]
  apply 	--application myapp
  unlock 	[ --application myapp ]
` + "```"
}

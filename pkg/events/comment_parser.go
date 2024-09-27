package events

import (
	"fmt"
	"io"
	"strings"
	"unicode"

	"github.com/corymurphy/argobot/pkg/command"
	"github.com/corymurphy/argobot/pkg/logging"
	"github.com/google/go-github/v53/github"
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
		AllowCommands: []command.Name{command.Help, command.Plan, command.Apply},
	}
}

func (c *CommentParser) Parse(event github.IssueCommentEvent) *CommentParseResult {
	if event.GetAction() != "created" {
		msg := "ignoring a non created event"
		c.Log.Info(msg)
		return &CommentParseResult{
			ImmediateResponse:  true,
			Ignore:             true,
			HasResponseComment: false,
		}
	}

	if !event.GetIssue().IsPullRequest() {
		msg := "ignoring a non pull request event"
		c.Log.Info(msg)
		return &CommentParseResult{
			ImmediateResponse:  true,
			Ignore:             true,
			HasResponseComment: false,
		}
	}

	author := event.GetComment().GetUser().GetLogin()
	isBot := strings.HasSuffix(author, "[bot]")
	if isBot {
		c.Log.Info("ignoring comment from a bot")
		return &CommentParseResult{
			ImmediateResponse:  true,
			Ignore:             true,
			HasResponseComment: false,
		}
	}

	body := event.GetComment().Body
	if len(*body) > CommandMaxLen {
		c.Log.Info("ignoring comment exceeding max length")
		return &CommentParseResult{
			ImmediateResponse:  true,
			Ignore:             true,
			HasResponseComment: false,
		}
	}

	isArgo := strings.HasPrefix(*body, CommandName)
	if !isArgo {
		c.Log.Debug("ignoring, comment is not an argobot command")
		return &CommentParseResult{
			ImmediateResponse:  true,
			Ignore:             true,
			HasResponseComment: false,
		}
	}

	*body = strings.ToLower(*body)

	if *body == CommandName || *body == fmt.Sprintf("%s help", CommandName) {
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

	args, err := shlex.Split(*body)
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
		flagSet = pflag.NewFlagSet(command.Apply.String(), pflag.ContinueOnError)
		flagSet.SetOutput(io.Discard)
		flagSet.StringVarP(&app, applicationFlagLong, applicationFlagShort, "", "ArgoCD application to run plan against")
	default:
		c.Log.Debug("failed to parse command %s", cmd)
		return &CommentParseResult{
			Ignore:            true,
			ImmediateResponse: true,
		}
	}

	flagSet.Parse(args[2:])

	if !c.isValidApplicationName(app) {
		c.Log.Debug("invalid application name %s", app)
		return &CommentParseResult{
			Ignore:            true,
			ImmediateResponse: true,
		}
	}

	return &CommentParseResult{
		Command: &CommentCommand{
			Name:        name,
			Application: app,
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

Commands:
  help 		Shows this message
  plan 		--application myapp
  apply 	--application myapp
` + "```"
}

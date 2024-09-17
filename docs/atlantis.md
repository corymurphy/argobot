# atlantis

How does atlantis accomplish this problem?

main - >  ServerCmd.init - returns cobra command with config
ServerCmd -> run
ServerCreator.NewServer - returns a server starter interface

- allowCommands, err := userConfig.ToAllowCommandNames()
- create github client
- create http router
- if github app, NewGithubAppTokenRotator
- create event parser
- create comment parser
- create pull updater (events.PullUpdater)
- events.CommentCommandRunner(events.NewPlanCommandRunner)
- events_controllers.VCSEventsController
- add all to a server object and return

- Start

package cmd

import (
	"log"
	"os/exec"
)

type ShellCommand struct {
	pathCache map[string]string
}

type Command interface {
	Run(string, ...string) ([]byte, error) // exe, args
}

func NewShellCommand() Command {
	return ShellCommand{
		pathCache: make(map[string]string),
	}
}

func NewCommand() Command {
	return ShellCommand{
		pathCache: make(map[string]string),
	}
}

func (c ShellCommand) Run(exe string, args ...string) ([]byte, error) {
	path := c.getExePath(exe)

	cmd := exec.Command(path, args...)

	if cmd.Err != nil {
		return nil, cmd.Err
	}

	return cmd.CombinedOutput()
}

func (c ShellCommand) getArgs(path string, args []string) []string {
	// cmd.Run requires the path to be the first element of the slice
	return append([]string{path}, args...)
}

func (c ShellCommand) getExePath(exe string) string {
	if path, hit := c.pathCache[exe]; hit { // cache hit
		return path
	}
	path, err := exec.LookPath(exe)
	if err != nil {
		log.Fatal(err)
	}

	c.pathCache[exe] = path

	return path
}

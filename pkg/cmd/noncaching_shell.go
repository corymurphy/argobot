package cmd

import "os/exec"

type NonCachingShell struct {
}

func (n NonCachingShell) Run(exe string, args ...string) ([]byte, error) {
	cmd := exec.Command(
		exe,
		args...,
	)

	if cmd.Err != nil {
		return nil, cmd.Err
	}

	return cmd.CombinedOutput()
}

func NewNonCachingShell() Command {
	return &NonCachingShell{}
}

package command

import "fmt"

type Name int

const (
	// Apply Name = iota
	Plan Name = iota
	Help Name = iota
)

func ParseCommandName(name string) (Name, error) {
	switch name {
	case "plan":
		return Plan, nil
	case "help":
		return Help, nil
	}
	return -1, fmt.Errorf("unknown command name: %s", name)
}

func (c Name) String() string {
	switch c {
	case Plan:
		return "plan"
	}
	return ""
}

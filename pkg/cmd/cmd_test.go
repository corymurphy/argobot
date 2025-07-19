package cmd

import (
	"testing"
)

func Test_Command_Echo_SuccessStatusCode(t *testing.T) {
	_, actual := NewCommand().Run("echo", "'hello world'")
	if actual != nil {
		t.Errorf("echo should have a zero exit code: %v", actual)
	}
}

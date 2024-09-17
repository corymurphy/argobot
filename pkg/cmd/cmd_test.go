package cmd

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Command_Echo_SuccessStatusCode(t *testing.T) {
	_, actual := NewCommand().Run("echo", "'hello world'")
	assert.Equal(t, nil, actual, "echo should have a zero exit code")
}

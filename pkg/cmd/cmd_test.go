package cmd

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Command_Echo_SuccessStatusCode(t *testing.T) {
	expected := 0
	actual, _ := NewCommand().Run("echo", []string{"'hello world'"})
	assert.Equal(t, expected, actual, "echo should have a zero exit code")
}

package assert

import (
	"fmt"
	"testing"
)

type TestLogger struct {
	t *testing.T
}

func NewTestLogger(t *testing.T) *TestLogger {
	return &TestLogger{
		t: t,
	}
}

func (t *TestLogger) log(format string, a ...interface{}) {
	err := fmt.Errorf(format, a...)
	t.t.Error(err)
}

func (t *TestLogger) Debug(format string, a ...interface{}) {
	t.log(format, a...)
}

func (t *TestLogger) Info(format string, a ...interface{}) {
	t.log(format, a...)
}

func (t *TestLogger) Warn(format string, a ...interface{}) {
	t.log(format, a...)
}

func (t *TestLogger) Err(format string, a ...interface{}) {
	t.log(format, a...)
}

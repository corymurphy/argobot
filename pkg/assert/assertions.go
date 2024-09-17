package assert

import (
	"testing"
)

func Equal(t *testing.T, expected interface{}, actual interface{}) {
	if expected != actual {
		t.Errorf("expected %s, got %s", expected, actual)
	}
}

func True(t *testing.T, value bool, message string) {
	if !value {
		t.Errorf("expected true for %s", message)
	}
}

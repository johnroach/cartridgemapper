package cmd

import (
	"testing"
)

func TestRootVersionCommand(t *testing.T) {
	if AppVersion != printVersion() {
		t.Errorf("printVersion() test failed.")
	}
}

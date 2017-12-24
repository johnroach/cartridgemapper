package cmd

import (
	"fmt"
	"testing"
)

func TestRootVersionCommand(t *testing.T) {
	fmt.Println("hullo")
	if AppVersion != printVersion() {
		t.Errorf("printVersion() test failed")
	}
}

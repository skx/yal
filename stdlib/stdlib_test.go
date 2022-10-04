package stdlib

import (
	"os"
	"testing"
)

func TestStdlib(t *testing.T) {
	x := Contents()

	var core []byte
	var mal []byte
	var err error

	core, err = os.ReadFile("stdlib.lisp")
	if err != nil {
		t.Fatalf("failed to read: %s", err)
	}

	mal, err = os.ReadFile("mal.lisp")
	if err != nil {
		t.Fatalf("failed to read: %s", err)
	}

	// Add one for the newline we add in the middle.
	total := len(core) + 1 + len(mal) + 1
	if total != len(x) {
		t.Fatalf("stdlib size mismatch")
	}

}

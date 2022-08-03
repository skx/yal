package stdlib

import (
	"os"
	"testing"
)

func TestStdlib(t *testing.T) {
	x := Contents()

	data, err := os.ReadFile("stdlib.lisp")
	if err != nil {
		t.Fatalf("failed to read: %s", err)
	}
	if len(data) != len(x) {
		t.Fatalf("stdlib size mismatch")
	}
	for i, b := range data {
		if x[i] != b {
			t.Fatalf("mismatch of contents at offset %d", i)
		}
	}
}

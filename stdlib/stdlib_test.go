package stdlib

import (
	"os"
	"testing"
)

// Test we can exclude all
func TestStdlibExcludeAll(t *testing.T) {

	// By default we get "stuff"
	x := Contents()

	if len(x) < 1 {
		t.Fatalf("Failed to get contents of stdlib")
	}

	// Excluding everything should return nothing
	os.Setenv("YAL_STDLIB_EXCLUDE_ALL", "yes")

	x = Contents()
	if len(x) != 0 {
		t.Fatalf("We expected no content, but got something, despite $YAL_STDLIB_EXCLUDE_ALL")
	}

	// restore
	os.Setenv("YAL_STDLIB_EXCLUDE_ALL", "")
}

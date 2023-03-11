package stdlib

import (
	"fmt"
	"os"
	"strings"
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
	fmt.Printf("%s\n", x)
	if len(x) != 0 {
		t.Fatalf("We expected no content, but got something, despite $YAL_STDLIB_EXCLUDE_ALL")
	}

	// restore
	os.Setenv("YAL_STDLIB_EXCLUDE_ALL", "")
}

// Test we can exclude time.lisp
func TestStdlibExcludeTime(t *testing.T) {

	// By default we get "stuff"
	x := Contents()

	if len(x) < 1 {
		t.Fatalf("Failed to get contents of stdlib")
	}

	// ensure we have "hms" function defined
	expected := "(set! time:hms"

	content := string(x)
	if !strings.Contains(content, expected) {
		t.Fatalf("failed to find expected function: %s", expected)
	}

	// Now exclude "time"
	os.Setenv("YAL_STDLIB_EXCLUDE", "time")

	// Re-read content
	x = Contents()
	if len(x) < 1 {
		t.Fatalf("Failed to get contents of stdlib")
	}

	content = string(x)
	if strings.Contains(content, expected) {
		t.Fatalf("we shouldn't find the excluded function, but we did: %s", expected)
	}

	// restore
	os.Setenv("YAL_STDLIB_EXCLUDE", "")
}

// Package stdlib contains a simple/small standard-library, which is written in lisp itself.
//
// By default our standard library is loaded prior to the execution of any user-supplied
// code, however parts of it can be selectively ignored, or the whole thing.
//
// If the environmental varialbe "YAL_STDLIB_EXCLUDE_ALL" contains non-empty content then
// all of our standard-library is disabled.
//
// Otherwise if YAL_STDLIB_EXCLUDE is set to a non-empty string it will be assumed to be
// a comma-separated list of filename substrings to exclude.
package stdlib

import (
	"embed" // embedded-resource magic
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

//go:embed stdlib/*.lisp
var stdlib embed.FS

// Contents returns the embedded contents of our Lisp standard-library.
//
// We embed "*.lisp" when we build our binary
func Contents() []byte {

	// Result
	result := []byte{}

	// We can allow disabling the stdlib.
	if os.Getenv("YAL_STDLIB_EXCLUDE_ALL") != "" {
		return result
	}

	// We might exclude only one/two files
	exclude := []string{}
	if os.Getenv("YAL_STDLIB_EXCLUDE") != "" {
		exclude = strings.Split(os.Getenv("YAL_STDLIB_EXCLUDE"), ",")
	}

	// Read the list of entries
	entries, err := stdlib.ReadDir("stdlib")
	if err != nil {
		fmt.Printf("Failed to read embedded resources; fatal error\n")
		return result
	}

	// For each entry
	for _, entry := range entries {

		// Get the filename
		fp := filepath.Join("stdlib", entry.Name())

		// Does this match an excluded value?
		skip := false

		for _, tmp := range exclude {
			if strings.Contains(fp, tmp) {
				skip = true
			}
		}

		if skip {
			fmt.Printf("Skipping %s\n", fp)
			continue
		}

		// Read the content
		data, err := stdlib.ReadFile(fp)
		if err != nil {
			fmt.Printf("Failed to read embedded resource - %s - fatal error %s\n", fp, err)
			return result
		}

		// Append to our result
		result = append(result, data...)
	}

	return result
}

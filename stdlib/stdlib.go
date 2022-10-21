// Package stdlib contains a simple & small standard-library,
// written in lisp itself.
//
// By default our standard library is loaded prior to the
// execution of any user-supplied code.
package stdlib

import (
	"embed" // embedded-resource magic
	"fmt"
	"path/filepath"
)

//go:embed *.lisp
var stdlib embed.FS

// Contents returns the embedded contents of our Lisp standard-library.
//
// We embed "*.lisp" when we build our binary.
func Contents() []byte {

	// Result
	result := []byte{}

	// Read the list of entries
	entries, err := stdlib.ReadDir(".")
	if err != nil {
		fmt.Printf("Failed to read embedded resources; fatal error\n")
		return result
	}

	// For each entry
	for _, entry := range entries {

		// Get the filename
		fp := filepath.Join(".", entry.Name())

		// Read the content
		data, err := stdlib.ReadFile(fp)
		if err != nil {
			fmt.Printf("Failed to read embedded resource - %s - fatal error\n", fp)
			return result
		}

		// Append to our result
		result = append(result, data...)
	}

	return result
}

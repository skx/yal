//go:build windows

package builtins

import (
	"fmt"
	"os"
)

// getGID should return the group of the file, from the extended information
// available after a stat, however on Windows platforms that doesn't work
// in the obvious way.
//
// Here we just return an error to make that apparent to the caller.
//
// This is in a separate file so that we use build-tags to build code
// appropriately.
func getGID(info os.FileInfo) (int, error) {

	return 0, fmt.Errorf("not found")
}

// getUID should return the owner of the file, from the extended information
// available after a stat, however on Windows platforms that doesn't work
// in the obvious way.
//
// Here we just return an error to make that apparent to the caller.
//
// This is in a separate file so that we use build-tags to build code
// appropriately.
func getUID(info os.FileInfo) (int, error) {

	return 0, fmt.Errorf("not found")
}

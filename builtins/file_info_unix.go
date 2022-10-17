//go:build !windows

package builtins

import (
	"fmt"
	"os"
	"syscall"
)

// getUID returns the owner of the file, from the extended information
// available after a stat.
//
// This is in a seperate file so that we can build upon Windows systems.
func getUID(info os.FileInfo) (int, error) {

	if stat, ok := info.Sys().(*syscall.Stat_t); ok {
		return int(stat.Uid), nil
	}
	return 0, fmt.Errorf("not found")
}

// getGID returns the group of the file, from the extended information
// available after a stat.
//
// This is in a seperate file so that we can build upon Windows systems.
func getGID(info os.FileInfo) (int, error) {

	if stat, ok := info.Sys().(*syscall.Stat_t); ok {
		return int(stat.Gid), nil
	}
	return 0, fmt.Errorf("not found")
}

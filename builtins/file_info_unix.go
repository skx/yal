//go:build !windows

package builtins

import (
	"os"
	"syscall"
)

// getGID returns the group of the file, from the extended information
// available after a stat - that is not portable to Windows though.
//
// This is in a separate file so that we use build-tags to build code
// appropriately.
func getGID(info os.FileInfo) (int, error) {

	stat, _ := info.Sys().(*syscall.Stat_t)
	return int(stat.Gid), nil
}

// getUID returns the owner of the file, from the extended information
// available after a stat - that is not portable to Windows though.
//
// This is in a separate file so that we use build-tags to build code
// appropriately.
func getUID(info os.FileInfo) (int, error) {

	stat, _ := info.Sys().(*syscall.Stat_t)
	return int(stat.Uid), nil
}

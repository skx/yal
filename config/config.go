// Package config provides an I/O abstraction for our interpreter,
// allowing it to be embedded and used in places where STDIN and STDOUT
// are not necessarily terminal-based.
//
// All input-reading uses the level of indirection provided here, and
// similarly output goes via the writer we hold here.
//
// This abstraction allows a host program to setup a different pair of
// streams prior to initializing the interpreter.
package config

import (
	"io"
	"os"
)

// Config is a holder for configuration which is used for interfacing
// the interpreter with the outside world.
type Config struct {

	// STDIN is an input-reader used for the (read) function, when
	// called with no arguments.
	STDIN io.Reader

	// STDOUT is the writer which is used for "(print)".
	STDOUT io.Writer
}

// New returns a new configuration object
func New() *Config {

	e := &Config{}
	return e
}

// DefaultIO returns a configuration which uses the default
// input and output streams - i.e. STDIN and STDOUT work as
// expected
func DefaultIO() *Config {
	e := New()

	// Setup default input/output streams
	e.STDIN = os.Stdin
	e.STDOUT = os.Stdout

	return e
}

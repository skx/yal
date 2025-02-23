// Package config provides an I/O abstraction for our interpreter,
// allowing it to be embedded and used in places where STDIN and STDOUT
// are not necessarily terminal-based.
//
// All input-reading uses the level of indirection provided here, and
// similarly output goes via the writer we hold here. There is also a
// STDERR stream which is used for (optional) debugging output by our
// main driver.
//
// The I/O abstraction allows a host program to setup different streams
// prior to initializing the interpreter.
package config

import (
	"io"
	"os"
)

// Config is a holder for configuration which is used for interfacing
// the interpreter with the outside world.
type Config struct {

	// STDERR is the writer for debug output, by default output
	// sent here is discarded.
	STDERR io.Writer

	// STDIN is an input-reader used for the (read) function, when
	// called with no arguments.
	STDIN io.Reader

	// STDOUT is the writer which is used for "(print)".
	STDOUT io.Writer
}

// New returns a new configuration object
func New() *Config {

	e := new(Config)
	return e
}

// DefaultIO returns a configuration which uses the default
// input and output streams - i.e. STDIN and STDOUT work as
// expected.
//
// The STDERR writer is configured to discard output by default.
func DefaultIO() *Config {
	e := New()

	// Setup useful input/output streams for default usage.
	e.STDIN = os.Stdin
	e.STDOUT = os.Stdout

	// STDERR is only used when debugging.
	//
	// So we can discard output here by default.
	e.STDERR = io.Discard

	return e
}

package primitive

import "fmt"

// Error holds an error message.
type Error string

// ToInterface converts this object to a golang value
func (e Error) ToInterface() any {
	return fmt.Errorf(string(e))
}

// ToString converts this object to a string.
func (e Error) ToString() string {
	return "ERROR{" + string(e) + "}"
}

// Type returns the type of this primitive object.
func (e Error) Type() string {
	return "error"
}

// ArityError is the error raised when a function, or special form,
// is invoked with the wrong number of arguments.
func ArityError() Error {
	return Error("ArityError - Unexpected argument count")
}

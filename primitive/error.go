package primitive

import "fmt"

// Error holds an error message.
type Error string

// ArityError is the error raised when a function, or special form,
// is invoked with the wrong number of arguments.
func ArityError() Error {
	return Error("ArityError - Unexpected argument count")
}

// TypeError is an error raised when a function is called with invalid
// typed argument
func TypeError(msg string) Error {
	return Error("TypeError - " + msg)
}

// IsSimpleType is used to denote whether this object
// is self-evaluating.
func (e Error) IsSimpleType() bool {
	return true
}

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

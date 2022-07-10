// Package primitive contains the definitions of our primitive types,
// which are "nil", "bool", "number", "string", and "list".
package primitive

// Primitive is the interface of all our types
type Primitive interface {

	// Convert this primitive to a string
	ToString() string

	// Return the type of this object
	Type() string
}

// IsNil tests whether an expression is nil.
func IsNil(e Primitive) bool {
	var n Nil
	return e == n
}

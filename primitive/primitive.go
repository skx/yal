// Package primitive contains the definitions of our primitive types,
// which are "nil", "bool", "number", "string", and "list".
package primitive

// Primitive is the interface of all our types
type Primitive interface {

	// IsSimpleType is used to denote whether this object
	// is self-evaluating.
	//
	// Simple types include strings, numbers, booleans, etc.
	//
	// However note that a list is NOT a simple type, as it
	// is used to denote a function-call.
	IsSimpleType() bool

	// ToString converts this primitive to a string representation.
	ToString() string

	// Type returns the type of this object.
	Type() string
}

// ToNative is an optional interface that some of our primitive
// types might choose to implement.
//
// If available this allows a YAL object to be converted to a
// suitable Golang equivalent type/value.
type ToNative interface {

	// ToInterface converts to a native golang type.
	ToInterface() interface{}
}

// IsNil tests whether an expression is nil.
func IsNil(e Primitive) bool {
	var n Nil
	return e == n
}

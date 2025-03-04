package primitive

import "fmt"

// Number type holds numbers.
type Number float64

// IsSimpleType is used to denote whether this object
// is self-evaluating.
func (n Number) IsSimpleType() bool {
	return true
}

// IsInt returns true if this number is an integer
func (n Number) IsInt() bool {
	return float64(n) == float64(int(n))
}

// ToInterface converts this object to a golang value
func (n Number) ToInterface() any {

	// int?
	if n.IsInt() {
		return int(n)
	}

	// float
	return float64(n)
}

// ToString converts this object to a string.
func (n Number) ToString() string {

	// Is this really an integer?
	if n.IsInt() {
		return fmt.Sprintf("%d", int(n))
	}

	return fmt.Sprintf("%f", n)
}

// Type returns the type of this primitive object.
func (n Number) Type() string {
	return "number"
}

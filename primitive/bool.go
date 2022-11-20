package primitive

// Bool is our wrapping of bool
type Bool bool

// IsSimpleType is used to denote whether this object
// is self-evaluating.
func (b Bool) IsSimpleType() bool {
	return true
}

// ToInterface converts this object to a golang value
func (b Bool) ToInterface() any {
	return bool(b)
}

// ToString converts this object to a string.
func (b Bool) ToString() string {
	if b {
		return "#t"
	}
	return "#f"
}

// Type returns the type of this primitive object.
func (b Bool) Type() string {
	return "boolean"
}

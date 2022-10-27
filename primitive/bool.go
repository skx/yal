package primitive

// Bool is our wrapping of bool
type Bool bool

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

package primitive

// Bool is our wrapping of bool
type Bool bool

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

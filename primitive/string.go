package primitive

// String holds a string value.
type String string

// IsSimpleType is used to denote whether this object
// is self-evaluating.
func (s String) IsSimpleType() bool {
	return true
}

// ToString converts this object to a string.
func (s String) ToString() string {
	return string(s)
}

// Type returns the type of this primitive object.
func (s String) Type() string {
	return "string"
}

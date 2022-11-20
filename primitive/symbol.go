package primitive

// Symbol is the type for our symbols.
type Symbol string

// IsSimpleType is used to denote whether this object
// is self-evaluating.
func (s Symbol) IsSimpleType() bool {
	return false
}

// ToInterface converts this object to a golang value
func (s Symbol) ToInterface() any {
	return s.ToString()
}

// ToString converts this object to a string.
func (s Symbol) ToString() string {
	return string(s)
}

// Type returns the type of this primitive object.
func (s Symbol) Type() string {
	return "symbol"
}

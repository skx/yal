package primitive

// Nil type holds the undefined value
type Nil struct{}

// ToString converts this object to a string.
func (n Nil) ToString() string {
	return "nil"
}

// Type returns the type of this primitive object.
func (n Nil) Type() string {
	return "nil"
}

package primitive

// Nil type holds the undefined value
type Nil struct{}

// ToInterface converts this object to a golang value
func (n Nil) ToInterface() any {
	return nil
}

// ToString converts this object to a string.
func (n Nil) ToString() string {
	return "nil"
}

// Type returns the type of this primitive object.
func (n Nil) Type() string {
	return "nil"
}

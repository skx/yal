package primitive

// Character holds a string value.
type Character string

// ToString converts this object to a string.
func (c Character) ToString() string {
	return string(c)
}

// Type returns the type of this primitive object.
func (c Character) Type() string {
	return "character"
}

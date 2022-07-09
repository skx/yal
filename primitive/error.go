package primitive

// Error holds an error message.
type Error string

// ToString converts this object to a string.
func (e Error) ToString() string {
	return "ERROR{" + string(e) + "}"
}

// Type returns the type of this primitive object.
func (e Error) Type() string {
	return "error"
}

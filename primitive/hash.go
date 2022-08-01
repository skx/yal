package primitive

// Hash holds a collection of other types, indexed by string
type Hash struct {
	Entries map[string]Primitive
}

// Get returns the value of a given index
func (h Hash) Get(key string) Primitive {
	return h.Entries[key]
}

// Set stores a value in the hash
func (h Hash) Set(key string, val Primitive) {
	h.Entries[key] = val
}

// ToString converts this object to a string.
func (h Hash) ToString() string {

	out := "{"
	for k, v := range h.Entries {
		out += " " + k + ":" + v.ToString()
	}
	out += "}"
	return out
}

// Type returns the type of this primitive object.
func (h Hash) Type() string {
	return "hash"
}

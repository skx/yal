package primitive

import "fmt"

// Hash holds a collection of other types, indexed by string
type Hash struct {

	// Entries contains the key/value pairs this object holds.
	Entries map[string]Primitive

	// StructType contains the name of this struct, if it is being
	// being used to implement a Struct, rather than a Hash
	StructType string
}

// Get returns the value of a given index
func (h Hash) Get(key string) Primitive {
	x, ok := h.Entries[key]
	if ok {
		return x
	}
	return Nil{}
}

// GetStruct returns the name of the structure this object contains, if any
func (h *Hash) GetStruct() string {
	return h.StructType
}

// IsSimpleType is used to denote whether this object
// is self-evaluating.
func (h Hash) IsSimpleType() bool {
	return true
}

// NewHash creates a new hash, and ensures that the storage-space
// is initialized.
func NewHash() Hash {
	h := Hash{}
	h.Entries = make(map[string]Primitive)
	return h
}

// Set stores a value in the hash
func (h Hash) Set(key string, val Primitive) {
	h.Entries[key] = val
}

// SetStruct marks this as a "struct" type instead of a "hash type",
// when queried by lisp
func (h *Hash) SetStruct(name string) {
	h.StructType = name
}

// ToString converts this object to a string.
func (h Hash) ToString() string {

	out := "{\n"
	for k, v := range h.Entries {
		out += "\t" + k + " => " + v.ToString() + "\n"
	}
	out += "}"
	return out
}

// Type returns the type of this primitive object.
func (h Hash) Type() string {
	if h.StructType == "" {
		return "hash"
	}
	return fmt.Sprintf("struct-%s", h.StructType)
}

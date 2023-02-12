package primitive

import "sort"

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
//
// Note that we sort the keys before returning the stringified object,
// which allows us to use the "eq" test on hashes with identical key/values,
// regardless of their ordering.
func (h Hash) ToString() string {

	// Output prefix.
	out := "{\n"

	// Get the keys in our hash.
	keys := []string{}

	for x := range h.Entries {
		keys = append(keys, x)
	}

	// Sort the list of keys
	sort.Strings(keys)

	// Now we can get a consistent ordering for our
	// hash keys/value pairs.
	for _, key := range keys {
		out += "\t" + key + " => " + h.Entries[key].ToString() + "\n"
	}

	// Terminate the string representation and return.
	out += "}"
	return out
}

// Type returns the type of this primitive object.
func (h Hash) Type() string {
	if h.StructType == "" {
		return "hash"
	}
	return h.StructType
}

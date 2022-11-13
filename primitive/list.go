package primitive

import "strings"

// List holds a collection of other types, including Lists.
type List []Primitive

// IsSimpleType is used to denote whether this object
// is self-evaluating.
func (l List) IsSimpleType() bool {
	return false
}

// ToString converts this object to a string.
func (l List) ToString() string {
	elemStrings := []string{}
	for _, e := range l {
		elemStrings = append(elemStrings, e.ToString())
	}
	return "(" + strings.Join(elemStrings, " ") + ")"
}

// Type returns the type of this primitive object.
func (l List) Type() string {
	return "list"
}

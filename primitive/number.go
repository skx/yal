package primitive

import (
	"strconv"
)

// Number type holds numbers
type Number float64

// ToString converts this object to a string.
func (n Number) ToString() string {
	return strconv.FormatFloat(float64(n), 'g', -1, 64)
}

// Type returns the type of this primitive object.
func (n Number) Type() string {
	return "number"
}

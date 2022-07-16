package primitive

import (
	"fmt"
)

// Number type holds numbers.
type Number float64

// ToString converts this object to a string.
func (n Number) ToString() string {

	// Is this really an integer?
	if float64(n) == float64(int(n)) {
		return fmt.Sprintf("%d", int(n))
	}

	return fmt.Sprintf("%f", n)
}

// Type returns the type of this primitive object.
func (n Number) Type() string {
	return "number"
}

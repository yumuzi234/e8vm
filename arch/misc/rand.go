package misc

import (
	"io"
)

// Rand provides a random number generator.
type Rand struct {
	r io.Reader
}

// NewRand creates a new random number generator.
func NewRand(seed int64) *Rand {
	panic("todo")
}

// Handle handles incoming request to generate a random number.
func (r *Rand) Handle(_ []byte) ([]byte, int32) {
	panic("todo")
}

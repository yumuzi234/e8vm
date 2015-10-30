package types

import (
	"fmt"
)

// Array is an array type of fixed size
type Array struct {
	T T
	N int32
}

// Size returns the total size of the array.
func (t *Array) Size() int32 {
	size := t.T.Size()
	if t.RegSizeAlign() {
		size = RegSizeAlignUp(size)
	}
	return size * t.N
}

// String returns "[N]T"
func (t *Array) String() string {
	return fmt.Sprintf("[%d]%s", t.N, t.T)
}

// RegSizeAlign inherits from the internal type.
func (t *Array) RegSizeAlign() bool {
	return t.T.RegSizeAlign()
}

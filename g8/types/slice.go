package types

import (
	"e8vm.io/e8vm/arch8"
)

// Slice is a slice type
type Slice struct{ T T }

// String returns "[]T"
func (t *Slice) String() string { return "[]" + t.T.String() }

// Size returns the size of the slice.
// It contains the start address of the slice,
// and the number of elements of the slice.
func (t *Slice) Size() int32 { return arch8.RegSize * 2 }

// SliceOf returns the internal type of the slice.
// If the type is not a slice, it returns nil.
func SliceOf(t T) T {
	st, ok := t.(*Slice)
	if !ok {
		return nil
	}
	return st.T
}

// RegSizeAlign returns true. A slice is always word aligned.
func (t *Slice) RegSizeAlign() bool { return true }

// String is a slice of int8's
var String T = &Slice{Int8}

package types

import (
	"e8vm.io/e8vm/arch8"
)

// RegSizeAlignUp aligns a size up to multiples of the register size.
func RegSizeAlignUp(size int32) int32 {
	mod := size % arch8.RegSize
	if mod == 0 {
		return size
	}
	return size + arch8.RegSize - mod
}

// Pointer is a pointer type
type Pointer struct{ T T } // a pointer type

// NewPointer returns a pointer type of a particular type
func NewPointer(t T) *Pointer { return &Pointer{t} }

// String returns "*T"
func (t *Pointer) String() string { return "*" + t.T.String() }

// Size returns the address length of the architecture.
func (t *Pointer) Size() int32 { return arch8.RegSize }

// RegSizeAlign returns true. Pointer is always word aligned.
func (t *Pointer) RegSizeAlign() bool { return true }

// PointerOf returns the internal type of the pointer type.
// If the type is not a pointer, it returns nil.
func PointerOf(t T) T {
	pt, ok := t.(*Pointer)
	if !ok {
		return nil
	}
	return pt.T
}

// IsPointer checks if the type is a pointer.
func IsPointer(t T) bool {
	_, ok := t.(*Pointer)
	return ok
}

package types

import (
	"shanhu.io/smlvm/arch"
	"shanhu.io/smlvm/syms"
)

// Interface is the type of a structure
type Interface struct {
	Syms *syms.Table
	name string
}

// NewInterface constructs a new struct type.
func NewInterface(name string) *Interface {
	syms := syms.NewTable()
	return &Interface{
		Syms: syms,
		name: name,
	}
}

// Size returns the overall size of the structure type
func (t *Interface) Size() int32 { return arch.RegSize * 2 }

// String returns the name of the structure type
func (t *Interface) String() string { return t.name }

// RegSizeAlign returns true when at least one field in the struct
// is word aligned.
func (t *Interface) RegSizeAlign() bool { return true }

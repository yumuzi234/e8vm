package types

import (
	"shanhu.io/smlvm/syms"
)

// Interface is the type of a structure
type Interface struct {
	Syms *syms.Table

	name         string
	size         int32
	regSizeAlign bool
}

// NewInterface constructs a new struct type.
func NewInterface(name string) *Interface {
	syms := syms.NewTable()
	return &Interface{
		Syms:         syms,
		name:         name,
		size:         8,
		regSizeAlign: true,
	}
}

// Size returns the overall size of the structure type
func (t *Interface) Size() int32 { return t.size }

// String returns the name of the structure type
func (t *Interface) String() string { return t.name }

// RegSizeAlign returns true when at least one field in the struct
// is word aligned.
func (t *Interface) RegSizeAlign() bool { return t.regSizeAlign }

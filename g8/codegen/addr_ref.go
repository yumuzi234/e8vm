package codegen

import (
	"fmt"
)

// AddrRef is an indirect reference of an address range.
// we need this to implement *pointer, array, and struct.
// that is all the other right-hand values that is more than
// just a simple temp variable.
type AddrRef struct {
	base         Ref
	offset       int32
	size         int32
	u8           bool
	regSizeAlign bool
}

// NewAddrRef creates an indirect reference to an address that is
// saved at base.
func NewAddrRef(base Ref, n, offset int32, u8, regSizeAlign bool) Ref {
	return newAddrRef(base, n, offset, u8, regSizeAlign)
}

func newAddrRef(base Ref, n, offset int32, u8, regSizeAlign bool) *AddrRef {
	return &AddrRef{
		base:         base,
		size:         n,
		offset:       offset,
		u8:           u8,
		regSizeAlign: regSizeAlign,
	}
}

func (r *AddrRef) String() string {
	if r.offset == 0 {
		return fmt.Sprintf("*%s", r.base.String())
	}
	return fmt.Sprintf("*(%s+%d)", r.base.String(), r.offset)
}

// Size returns the size of the address reference
func (r *AddrRef) Size() int32 { return r.size }

// RegSizeAlign tells if the refernece is register size aligned
func (r *AddrRef) RegSizeAlign() bool { return r.regSizeAlign }

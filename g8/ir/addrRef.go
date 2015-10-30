package ir

import (
	"fmt"
)

// addrRef is an indirect reference of an address range.
// we need this to implement *pointer, array, and struct.
// that is all the other right-hand values that is more than
// just a simple temp variable.
type addrRef struct {
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

func newAddrRef(base Ref, n, offset int32, u8, regSizeAlign bool) *addrRef {
	ret := new(addrRef)
	ret.base = base
	ret.size = n
	ret.offset = offset
	ret.u8 = u8
	ret.regSizeAlign = regSizeAlign

	return ret
}

func (r *addrRef) String() string {
	if r.offset == 0 {
		return fmt.Sprintf("*%s", r.base.String())
	}
	return fmt.Sprintf("*(%s+%d)", r.base.String(), r.offset)
}

func (r *addrRef) Size() int32 { return r.size }

func (r *addrRef) RegSizeAlign() bool { return r.regSizeAlign }

package ir

// varRef is a variable on stack, it is always word aligned
type varRef struct {
	name string // not unique, just for debugging

	// the offset relative to SP
	// before SP shift, the variable is saved at [SP-offset]
	// after SP shift, the variable is saved at [SP+framesize-offset]
	offset       int32
	size         int32
	u8           bool
	regSizeAlign bool

	// reg is the register allocated
	// valid values are in range [1, 4] for normal values
	// and also ret register is 6
	viaReg uint32
}

func newVar(n int32, name string, u8, regSizeAlign bool) *varRef {
	ret := new(varRef)
	ret.name = name
	ret.size = n
	ret.u8 = u8
	ret.regSizeAlign = regSizeAlign

	return ret
}

func (v *varRef) String() string {
	if v.name != "" {
		return v.name
	}
	return "<?>"
}

func (v *varRef) Size() int32 { return v.size }

func (v *varRef) RegSizeAlign() bool { return v.regSizeAlign }

func (v *varRef) canViaReg() bool {
	return v.size == 1 || v.size == 4
}

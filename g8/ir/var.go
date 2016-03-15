package ir

// Var is a variable on stack, it is always word aligned
type Var struct {
	name string // not unique, just for debugging

	// the offset relative to SP
	// before SP shift, the variable is saved at [SP-offset]
	// after SP shift, the variable is saved at [SP+framesize-offset]
	Offset       int32
	size         int32
	U8           bool
	regSizeAlign bool

	// reg is the register allocated
	// valid values are in range [1, 4] for normal values
	// and also ret register is 6
	viaReg uint32
}

// NewVar creates a new variable.
func NewVar(n int32, name string, u8, regSizeAlign bool) *Var {
	ret := new(Var)
	ret.name = name
	ret.size = n
	ret.U8 = u8
	ret.regSizeAlign = regSizeAlign

	return ret
}

func (v *Var) String() string {
	if v.name != "" {
		return v.name
	}
	return "<?>"
}

// Size returns the size of the variable
func (v *Var) Size() int32 { return v.size }

// RegSizeAlign tells if the variable is register size aligned.
func (v *Var) RegSizeAlign() bool { return v.regSizeAlign }

// CanViaReg tells if the variables can be saved and loaded via
// a register.
func (v *Var) CanViaReg() bool {
	return v.size == 1 || v.size == 4
}

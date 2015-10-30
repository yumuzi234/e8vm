package dasm8

// NewLine disassembles one instrcution at address addr.
func NewLine(addr uint32, in uint32) *Line {
	if (in >> 31) == 1 {
		return instJmp(addr, in)
	}

	op := (in >> 24) & 0xff
	switch {
	case op == 0:
		return instReg(addr, in)
	case op < 32:
		return instImm(addr, in)
	case op < 64:
		return instBr(addr, in)
	case op < 128:
		return instSys(addr, in)
	default:
		panic("bug")
	}
}

// LineStr returns the assembly string of a instruction as if
// if is at address 0
func LineStr(in uint32) string {
	return NewLine(0, in).Str
}

package inst

// Jmp makes a jump instruction
func Jmp(op uint32, im int32) uint32 {
	ret := (op & 0x3) << 30
	ret |= uint32(im) & 0x3fffffff
	return ret
}

// Sys makes a system instruction
func Sys(op, reg uint32) uint32 {
	return ((op & 0xff) << 24) | ((reg & 0x7) << 21)
}

// Imm composes an immediate based instruction
func Imm(op, d, s, im uint32) uint32 {
	ret := (op & 0xff) << 24
	ret |= (d & 0x7) << 21
	ret |= (s & 0x7) << 18
	ret |= (im & 0xffff)
	return ret
}

// Reg composes a register based instruction
func Reg(fn, d, s1, s2, sh, isFloat uint32) uint32 {
	ret := uint32(0)
	ret |= (d & 0x7) << 21
	ret |= (s1 & 0x7) << 18
	ret |= (s2 & 0x7) << 15
	ret |= (sh & 0x1f) << 10
	ret |= (isFloat & 0x1) << 8
	ret |= fn & 0xff
	return ret
}

// Br compose a branch instruction
func Br(op, s1, s2 uint32, im int32) uint32 {
	ret := (op & 0xff) << 24
	ret |= (s1 & 0x7) << 21
	ret |= (s2 & 0x7) << 18
	ret |= uint32(im) & 0x3ffff
	return ret
}

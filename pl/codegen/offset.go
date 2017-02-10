package codegen

func addOffsetHigh(b *Block, reg uint32, offset int32) {
	u := (uint32)(offset)
	sign := (u >> 15) & 0x1
	hi := ((u >> 16) + sign) & 0xffff
	if hi > 0 {
		b.inst(asm.addui(reg, reg, hi))
	}
}

package codegen

func addOffsetHigh(b *Block, reg uint32, offset int32) {
	if offset >= 0 {
		hi := uint32(offset) >> 16
		if hi > 0 {
			b.inst(asm.addui(reg, reg, hi))
		}
		return
	}

	// addi will add a 0xffff on the high 16-bit
	// we need to add 1 extra to cancel that out.
	hi := ((uint32(offset) >> 16) + 1) & 0xffff
	if hi > 0 {
		b.inst(asm.addui(reg, reg, hi))
	}
}

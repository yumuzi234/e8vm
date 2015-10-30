package ir

func genBlock(g *gener, b *Block) {
	for _, op := range b.ops {
		// printOp(os.Stderr, op)
		genOp(g, b, op)
	}

	// printJump(os.Stderr, b.jump)

	if b.jump == nil {
		/* do nothing */
	} else if b.jump.typ == jmpAlways {
		b.jumpInst = b.inst(asm.j(0))
	} else if b.jump.typ == jmpIfNot {
		loadRef(b, _4, b.jump.cond)
		b.jumpInst = b.inst(asm.beq(_0, _4, 0))
	} else if b.jump.typ == jmpIf {
		loadRef(b, _4, b.jump.cond)
		b.jumpInst = b.inst(asm.bne(_0, _4, 0))
	} else {
		panic("unknown jump")
	}
}

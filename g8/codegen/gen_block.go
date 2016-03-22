package codegen

func genBlock(g *gener, b *Block) {
	for _, op := range b.ops {
		genOp(g, b, op)
	}

	if b.jump == nil {
		/* do nothing */
	} else if b.jump.typ == jmpAlways {
		b.jumpInst = b.inst(asm.j(0))
	} else if b.jump.typ == jmpIfNot {
		loadRef(b, _r4, b.jump.cond)
		b.jumpInst = b.inst(asm.beq(_r0, _r4, 0))
	} else if b.jump.typ == jmpIf {
		loadRef(b, _r4, b.jump.cond)
		b.jumpInst = b.inst(asm.bne(_r0, _r4, 0))
	} else {
		panic("unknown jump")
	}
}

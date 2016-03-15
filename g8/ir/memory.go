package ir

import (
	"fmt"
)

func loadAddr(b *Block, reg uint32, r Ref) {
	switch r := r.(type) {
	case *Var:
		b.inst(asm.addi(reg, _sp, *b.frameSize-r.Offset))
	case *addrRef:
		loadRef(b, reg, r.base)
		if r.offset != 0 {
			b.inst(asm.addi(reg, reg, r.offset))
		}
	case *HeapSym:
		loadSym(b, reg, r.pkg, r.name)
	case *testList:
		loadSym(b, reg, r.pkg, r.name)
	default:
		panic(fmt.Errorf("load addr of %T", r))
	}
}

func loadArgAddr(b *Block, reg uint32, r *Var) {
	b.inst(asm.addi(reg, _sp, -r.Offset))
}

package ir

import (
	"fmt"

	"e8vm.io/e8vm/link8"
)

func isFuncPointer(r Ref) bool {
	if _, ok := r.(*Func); ok {
		return true
	}
	if _, ok := r.(*FuncSym); ok {
		return true
	}
	return false
}

func copyRef(g *gener, b *Block, dest, src Ref, isArg bool) {
	loadDestAddr := func(r uint32) {
		if !isArg {
			loadAddr(b, r, dest)
		} else {
			loadArgAddr(b, r, dest.(*varRef))
		}
	}

	size := dest.Size()
	switch {
	case size != src.Size():
		e := fmt.Errorf("copyRef src(%T)=%d dest(%T)=%d",
			src, src.Size(), dest, size)
		panic(e)
	case size < 0:
		panic("negative size for copyRef")
	case size == 0:
		return
	case canViaReg(dest) && canViaReg(src):
		loadRef(b, _4, src)
		if !isArg {
			saveRef(b, _4, dest, _1)
		} else {
			saveArg(b, _4, dest.(*varRef))
		}
	case isFuncPointer(src):
		loadRef(b, _4, src)
		loadDestAddr(_1)
		b.inst(asm.sw(_4, _1, 0))
		b.inst(asm.sw(_0, _1, regSize))
	default:
		if !isArg {
			loadAddr(b, _1, dest)
		} else {
			loadArgAddr(b, _1, dest.(*varRef))
		}
		loadAddr(b, _2, src)
		loadUint32(b, _3, uint32(size))

		jal := b.inst(asm.jal(0))
		f := g.memCopy
		jal.sym = &linkSym{link8.FillLink, f.pkg, f.sym}
	}
}

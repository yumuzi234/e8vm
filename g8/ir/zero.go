package ir

import (
	"fmt"

	"e8vm.io/e8vm/link8"
)

func zeroAddr(g *gener, b *Block, reg uint32, size int32, regSizeAlign bool) {
	switch {
	case size < 4:
		for i := int32(0); i < size; i++ {
			b.inst(asm.sb(_0, _1, i))
		}
	case size == 4 && regSizeAlign:
		b.inst(asm.sw(_0, _1, 0))
	case size == 8 && regSizeAlign:
		b.inst(asm.sw(_0, _1, 0))
		b.inst(asm.sw(_0, _1, 4))
	default:
		loadUint32(b, _2, uint32(size))
		jal := b.inst(asm.jal(0))
		f := g.memClear
		jal.sym = &linkSym{link8.FillLink, f.pkg, f.name}
	}
}

func zeroRef(g *gener, b *Block, r Ref) {
	switch r := r.(type) {
	case *varRef:
		if r.size < 0 {
			panic("invalid varRef size")
		}

		switch r.size {
		case 0: // do nothing
		case 1, regSize:
			saveVar(b, 0, r)
		default:
			loadAddr(b, _1, r)
			loadUint32(b, _2, uint32(r.size))

			jal := b.inst(asm.jal(0))
			f := g.memClear
			jal.sym = &linkSym{link8.FillLink, f.pkg, f.name}
		}
	case *addrRef:
		if r.size == 0 {
			return
		}
		loadAddr(b, _1, r)
		zeroAddr(g, b, _1, r.size, r.regSizeAlign)
	case *HeapSym:
		if r.size == 0 {
			return
		}
		loadAddr(b, _1, r)
		zeroAddr(g, b, _1, r.size, true)
	case *number:
		panic("number are read only")
	default:
		panic(fmt.Errorf("not implemented: %T", r))
	}
}

// CanBeZero checks if a reference could be zero value.
func CanBeZero(r Ref) bool {
	switch r := r.(type) {
	case *number:
		return r.v == 0
	case *FuncSym:
		return false
	case *Func:
		return false
	}
	return true
}

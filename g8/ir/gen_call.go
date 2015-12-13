package ir

import (
	"fmt"

	"e8vm.io/e8vm/link8"
)

func callSig(op *callOp) *FuncSig {
	switch f := op.f.(type) {
	case *FuncSym:
		return f.sig
	case *Func:
		return f.sig
	case *FuncPtr:
		return f.sig
	default:
		panic(fmt.Errorf("non-callable: %T", f))
	}
}

func genCallOp(g *gener, b *Block, op *callOp) {
	sig := callSig(op)

	// load the args
	// first copy the ones that send in on the stack
	for i, arg := range sig.args {
		if arg.viaReg == 0 {
			copyRef(g, b, arg, op.args[i], true)
		}
	}
	// then set the ones loaded by register
	for i, arg := range sig.args {
		if arg.viaReg > 0 {
			loadRef(b, arg.viaReg, op.args[i])
		}
	}

	// do the function call
	switch f := op.f.(type) {
	case *FuncSym:
		jal := b.inst(asm.jal(0))
		jal.sym = &linkSym{link8.FillLink, f.pkg, f.name}
	case *Func:
		jal := b.inst(asm.jal(0))
		jal.sym = &linkSym{link8.FillLink, f.pkg, f.name}
	case *FuncPtr:
		// function pointer, set PC manually
		loadRef(b, _4, f.Ref)
		b.inst(asm.addi(_ret, _pc, 4))
		b.inst(asm.ori(_pc, _4, 0))
	default:
		panic("bug")
	}

	// save the returns
	// first save the ones returned via register
	for i, ret := range sig.rets {
		if ret.viaReg > 0 {
			saveRef(b, ret.viaReg, op.dest[i], _4)
		}
	}
	// then copy the ones stored on the stack
	for i, ret := range sig.rets {
		if ret.viaReg == 0 {
			loadAddr(b, _1, op.dest[i])
			loadArgAddr(b, _2, ret)
			loadUint32(b, _3, uint32(ret.Size()))
			jal := b.inst(asm.jal(0))
			f := g.memCopy
			jal.sym = &linkSym{link8.FillLink, f.pkg, f.name}
		}
	}
}

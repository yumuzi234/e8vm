package codegen

import (
	"fmt"

	"e8vm.io/e8vm/link8"
)

func callSig(op *CallOp) *FuncSig {
	switch f := op.F.(type) {
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

func genCallOp(g *gener, b *Block, op *CallOp) {
	sig := callSig(op)

	// load the args
	// first copy the ones that send in on the stack
	for i, arg := range sig.args {
		if arg.ViaReg == 0 {
			copyRef(g, b, arg, op.Args[i], true)
		}
	}
	// then set the ones loaded by register
	for i, arg := range sig.args {
		if arg.ViaReg > 0 {
			loadRef(b, arg.ViaReg, op.Args[i])
		}
	}

	// do the function call
	switch f := op.F.(type) {
	case *FuncSym:
		jal := b.inst(asm.jal(0))
		jal.sym = &linkSym{link8.FillLink, f.pkg, f.name}
	case *Func:
		jal := b.inst(asm.jal(0))
		jal.sym = &linkSym{link8.FillLink, f.pkg, f.name}
	case *FuncPtr:
		// function pointer, set PC manually
		loadRef(b, _r4, f.Ref)
		b.inst(asm.addi(_ret, _pc, 4))
		b.inst(asm.ori(_pc, _r4, 0))
	default:
		panic("bug")
	}

	// save the returns
	// first save the ones returned via register
	for i, ret := range sig.rets {
		if ret.ViaReg > 0 {
			saveRef(b, ret.ViaReg, op.Dest[i], _r4)
		}
	}
	// then copy the ones stored on the stack
	for i, ret := range sig.rets {
		if ret.ViaReg == 0 {
			loadAddr(b, _r1, op.Dest[i])
			loadArgAddr(b, _r2, ret)
			loadUint32(b, _r3, uint32(ret.Size()))
			jal := b.inst(asm.jal(0))
			f := g.memCopy
			jal.sym = &linkSym{link8.FillLink, f.pkg, f.name}
		}
	}
}

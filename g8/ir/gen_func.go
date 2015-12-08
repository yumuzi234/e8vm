package ir

import (
	"e8vm.io/e8vm/arch8"
)

/**
## E8VM calling convention

- r0, always keep it as zero, never touch it
- r1, the first arg or return, if not used, should keep the value
- r2, the second arg or return
- r3, the third arg or return
- r4, free form temp
- sp, stack pointer
- ret, return address
- pc, the program counter

other args are pushed on the stack

## Function Prologue
- push ret to the stack
- push r1-r3 to the stack, for archive

## Function Epilogue
- save ret values to the stack

**/

// pushVar allocates a frame slot for the local var
func pushVar(f *Func, vars ...*varRef) {
	for _, v := range vars {
		f.frameSize += v.size
		if v.regSizeAlign {
			f.frameSize = alignUp(f.frameSize, regSize)
		}
		v.offset = f.frameSize
	}
}

func layoutLocals(f *Func) {
	for i, used := range f.sig.argRegUsed {
		if used || i == 0 {
			continue
		}

		// the caller is not using this reg for sending
		// the argument, the callee hence needs to
		// save this register
		v := newVar(regSize, "", false, true)
		v.viaReg = uint32(i)
		f.savedRegs = append(f.savedRegs, v)
	}

	// layout the variables in the function
	f.frameSize = f.sig.frameSize
	f.retAddr = newVar(regSize, "", false, true)
	f.retAddr.viaReg = arch8.RET // the return address

	// if all args and rets are via register
	// then f.retAddr.offset should be -4 (offset=4), it is the nearest to SP
	pushVar(f, f.retAddr)
	pushVar(f, f.sig.regArgs...)
	pushVar(f, f.sig.regRets...)
	pushVar(f, f.savedRegs...)
	pushVar(f, f.locals...)

	// pad up
	f.frameSize = alignUp(f.frameSize, regSize)
}

func makeMainPrologue(f *Func) {
	b := f.prologue
	b.inst(asm.addi(_sp, _sp, -f.frameSize))
}

func makeMainEpilogue(f *Func) {
	b := f.epilogue
	b.inst(asm.addi(_sp, _sp, f.frameSize))
	b.inst(asm.halt())
}

func makePrologue(g *gener, f *Func) {
	b := f.prologue

	saveRetAddr(b, f.retAddr)
	// move the sp
	b.inst(asm.addi(_sp, _sp, -f.frameSize))

	for _, v := range f.sig.args {
		if v.viaReg == 0 {
			continue // skip args not sent in via register
		}
		saveVar(b, v.viaReg, v)
	}

	// this is for restoreing the registers
	for _, v := range f.savedRegs {
		saveVar(b, v.viaReg, v)
	}
}

func makeEpilogue(g *gener, f *Func) {
	b := f.epilogue

	for _, v := range f.savedRegs {
		loadVar(b, v.viaReg, v) // restoring the registers
	}

	for _, v := range f.sig.rets {
		if v.viaReg == 0 {
			continue
		}
		loadVar(b, v.viaReg, v)
	}

	b.inst(asm.addi(_sp, _sp, f.frameSize))
	// back to the caller
	loadRetAddr(b, f.retAddr)
}

func genFunc(g *gener, f *Func) {
	layoutLocals(f)
	if f.frameSize > arch8.PageSize {
		g.Errorf(f.pos, "stack too large in function %s", f.name)
		return
	}

	if f.isMain {
		makeMainPrologue(f)
		makeMainEpilogue(f)
	} else {
		makePrologue(g, f)
		makeEpilogue(g, f)
	}

	for b := f.prologue.next; b != f.epilogue; b = b.next {
		genBlock(g, b)
	}

	// TODO: check ranges
	// now setup the jumps
	ninst := int32(0)
	for b := f.prologue; b != nil; b = b.next {
		b.instStart = ninst
		ninst += int32(len(b.insts))
		b.instEnd = ninst
	}

	for b := f.prologue; b != nil; b = b.next {
		if b.jumpInst == nil {
			continue
		}

		delta := b.jump.to.instStart - b.instEnd
		switch b.jump.typ {
		case jmpAlways:
			b.jumpInst.inst = asm.j(delta)
		case jmpIfNot:
			// TODO: check in jump range
			b.jumpInst.inst = asm.beq(_0, _4, delta)
		case jmpIf:
			b.jumpInst.inst = asm.bne(_0, _4, delta)
		default:
			panic("bug")
		}
	}
}

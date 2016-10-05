package codegen

import (
	"math"

	"shanhu.io/smlvm/arch"
)

/**
## Calling convention

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
func pushVar(f *Func, vars ...*Var) {
	for _, v := range vars {
		f.frameSize += v.size
		if v.regSizeAlign {
			f.frameSize = alignUp(f.frameSize, regSize)
		}
		v.Offset = f.frameSize
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
		v := NewVar(regSize, "", false, true)
		v.ViaReg = uint32(i)
		f.savedRegs = append(f.savedRegs, v)
	}

	// layout the variables in the function
	f.frameSize = f.sig.frameSize
	f.retAddr = NewVar(regSize, "", false, true)
	f.retAddr.ViaReg = arch.RET // the return address

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

func makePrologue(g *gener, f *Func) {
	b := f.prologue

	saveRetAddr(b, f.retAddr)
	// move the sp
	b.inst(asm.addi(_sp, _sp, -f.frameSize))

	for _, v := range f.sig.args {
		if v.ViaReg == 0 {
			continue // skip args not sent in via register
		}
		saveVar(b, v.ViaReg, v)
	}

	// this is for restoreing the registers
	for _, v := range f.savedRegs {
		saveVar(b, v.ViaReg, v)
	}
}

func makeEpilogue(g *gener, f *Func) {
	b := f.epilogue

	for _, v := range f.savedRegs {
		loadVar(b, v.ViaReg, v) // restoring the registers
	}

	for _, v := range f.sig.rets {
		if v.ViaReg == 0 {
			continue
		}
		loadVar(b, v.ViaReg, v)
	}

	b.inst(asm.addi(_sp, _sp, f.frameSize))
	// back to the caller
	loadRetAddr(b, f.retAddr)
}

func genFunc(g *gener, f *Func) {
	layoutLocals(f)
	if f.frameSize > arch.PageSize {
		g.Errorf(f.pos, "stack too large in function %s", f.name)
		return
	}

	makePrologue(g, f)
	makeEpilogue(g, f)

	for b := f.prologue.next; b != f.epilogue; b = b.next {
		genBlock(g, b)
	}

	// now setup the jumps
	ninst := int32(0)
	for b := f.prologue; b != nil; b = b.next {
		n := len(b.insts)
		if n > math.MaxInt32/4 {
			g.Errorf(f.pos, "basic block too large in function %s", f.name)
			return
		}

		b.instStart = ninst
		ninst += int32(n)
		b.instEnd = ninst

		if ninst > math.MaxInt32/4 {
			g.Errorf(f.pos, "function %s too large", f.name)
			return
		}
	}

	deltaCheck := func(delta int32) {
		if delta > 0x1ffff || delta < -0x20000 {
			g.Errorf(f.pos, "branch out of range in function: %s", f.name)
		}
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
			deltaCheck(delta)
			b.jumpInst.inst = asm.beq(_r0, _r4, delta)
		case jmpIf:
			deltaCheck(delta)
			b.jumpInst.inst = asm.bne(_r0, _r4, delta)
		default:
			panic("bug")
		}
	}
}

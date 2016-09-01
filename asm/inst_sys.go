package asm

import (
	"e8vm.io/e8vm/arch"
	asminst "e8vm.io/e8vm/asm/inst"
	"e8vm.io/e8vm/lexing"
)

var (
	// op
	opSysMap = map[string]uint32{
		"halt":    arch.HALT,
		"syscall": arch.SYSCALL,
		"iret":    arch.IRET,
		"sleep":   arch.SLEEP,
	}

	// op reg
	opSys1Map = map[string]uint32{
		"jruser": arch.JRUSER,
		"vtable": arch.VTABLE,
	}

	// op reg reg
	opSys2Map = map[string]uint32{
		"sysinfo": arch.SYSINFO,
	}
)

func makeInstSys(op, reg1, reg2 uint32) *inst {
	return &inst{inst: asminst.Sys(op, reg1, reg2)}
}

func resolveInstSys(p lexing.Logger, ops []*lexing.Token) (*inst, bool) {
	op0 := ops[0]
	opName := op0.Lit
	args := ops[1:]
	var op, reg1, reg2 uint32

	argCount := func(n int) bool {
		if !argCount(p, ops, n) {
			return false
		}

		if n >= 1 {
			reg1 = resolveReg(p, args[0])
		}
		if n >= 2 {
			reg2 = resolveReg(p, args[1])
		}
		return true
	}

	var found bool
	if op, found = opSysMap[opName]; found {
		// op
		argCount(0)
	} else if op, found = opSys1Map[opName]; found {
		// op reg
		argCount(1)
	} else if op, found = opSys2Map[opName]; found {
		argCount(2)
	} else {
		return nil, false
	}

	return makeInstSys(op, reg1, reg2), true
}

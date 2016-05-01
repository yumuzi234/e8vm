package asm8

import (
	"e8vm.io/e8vm/arch8"
	asminst "e8vm.io/e8vm/asm8/inst"
	"e8vm.io/e8vm/lex8"
)

var (
	// op
	opSysMap = map[string]uint32{
		"halt":    arch8.HALT,
		"syscall": arch8.SYSCALL,
		"iret":    arch8.IRET,
	}

	// op reg
	opSys1Map = map[string]uint32{
		"jruser":  arch8.JRUSER,
		"vtable":  arch8.VTABLE,
		"sysinfo": arch8.SYSINFO,
	}
)

func makeInstSys(op, reg uint32) *inst {
	return &inst{inst: asminst.Sys(op, reg)}
}

func resolveInstSys(p lex8.Logger, ops []*lex8.Token) (*inst, bool) {
	op0 := ops[0]
	opName := op0.Lit
	args := ops[1:]
	var op, reg uint32

	argCount := func(n int) bool {
		if !argCount(p, ops, n) {
			return false
		}

		if n >= 1 {
			reg = resolveReg(p, args[0])
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
	} else {
		return nil, false
	}

	return makeInstSys(op, reg), true
}

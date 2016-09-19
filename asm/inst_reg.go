package asm

import (
	"strconv"

	"shanhu.io/smlvm/arch"
	asminst "shanhu.io/smlvm/asm/inst"
	"shanhu.io/smlvm/lexing"
)

var (
	// op reg reg shift
	opShiftMap = map[string]uint32{
		"sll": arch.SLL,
		"srl": arch.SRL,
		"sra": arch.SRA,
	}

	// op reg reg reg
	opReg3Map = map[string]uint32{
		"sllv": arch.SLLV,
		"srlv": arch.SRLV,
		"srla": arch.SRLA,
		"add":  arch.ADD,
		"sub":  arch.SUB,
		"and":  arch.AND,
		"or":   arch.OR,
		"xor":  arch.XOR,
		"nor":  arch.NOR,
		"slt":  arch.SLT,
		"sltu": arch.SLTU,
		"mul":  arch.MUL,
		"mulu": arch.MULU,
		"div":  arch.DIV,
		"divu": arch.DIVU,
		"mod":  arch.MOD,
		"modu": arch.MODU,
	}

	// op reg reg
	opReg2Map = map[string]uint32{
		"mov": arch.SLL,
	}

	// op reg reg reg
	opFloatMap = map[string]uint32{
		"fadd": arch.FADD,
		"fsub": arch.FSUB,
		"fmul": arch.FMUL,
		"fdiv": arch.FDIV,
		"fint": arch.FINT,
	}
)

func parseShift(p lexing.Logger, op *lexing.Token) uint32 {
	ret, e := strconv.ParseUint(op.Lit, 0, 32)
	if e != nil {
		p.Errorf(op.Pos, "invalid shift %q: %s", op.Lit, e)
		return 0
	}

	if (ret & 0x1f) != ret {
		p.Errorf(op.Pos, "shift too large: %s", op.Lit)
		return 0
	}

	return uint32(ret)
}

func makeInstReg(fn, d, s1, s2, sh, isFloat uint32) *inst {
	ret := asminst.Reg(fn, d, s1, s2, sh, isFloat)
	return &inst{inst: ret}
}

func resolveInstReg(log lexing.Logger, ops []*lexing.Token) (*inst, bool) {
	op0 := ops[0]
	opName := op0.Lit
	args := ops[1:]

	// common args
	var fn, d, s1, s2, sh, isFloat uint32

	argCount := func(n int) bool {
		if !argCount(log, ops, n) {
			return false
		}
		if n >= 2 {
			d = resolveReg(log, args[0])
			s1 = resolveReg(log, args[1])
		}
		return true
	}

	var found bool
	if opName == "panic" {
		// panic
	} else if fn, found = opShiftMap[opName]; found {
		// op reg reg shift
		if argCount(3) {
			sh = parseShift(log, args[2])
		}
	} else if fn, found = opReg3Map[opName]; found {
		// op reg reg reg
		if argCount(3) {
			s2 = resolveReg(log, args[2])
		}
	} else if fn, found = opReg2Map[opName]; found {
		// op reg reg
		argCount(2)
	} else if fn, found = opFloatMap[opName]; found {
		// op reg reg reg
		if argCount(3) {
			s2 = resolveReg(log, args[2])
		}
		isFloat = 1
	} else {
		return nil, false
	}

	return makeInstReg(fn, d, s1, s2, sh, isFloat), true
}

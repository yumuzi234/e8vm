package asm

import (
	"e8vm.io/e8vm/arch"
	asminst "e8vm.io/e8vm/asm/inst"
	"e8vm.io/e8vm/lexing"
)

var (
	// op reg reg label
	opBrMap = map[string]uint32{
		"bne": arch.BNE,
		"beq": arch.BEQ,
	}
)

func makeInstBr(op, s1, s2 uint32) *inst {
	ret := asminst.Br(op, s1, s2, 0)
	return &inst{inst: ret}
}

func resolveInstBr(p lexing.Logger, ops []*lexing.Token) (*inst, bool) {
	op0 := ops[0]
	opName := op0.Lit
	args := ops[1:]

	var (
		op, s1, s2 uint32
		lab        string
		symTok     *lexing.Token

		found bool
	)

	if op, found = opBrMap[opName]; found {
		// op reg reg label
		if argCount(p, ops, 3) {
			s1 = resolveReg(p, args[0])
			s2 = resolveReg(p, args[1])
			symTok = args[2]
			if checkLabel(p, symTok) {
				lab = symTok.Lit
			} else {
				p.Errorf(symTok.Pos, "expects a label for %s", opName)
			}
		}
	} else {
		return nil, false
	}

	ret := makeInstBr(op, s1, s2)
	ret.sym = lab
	ret.fill = fillLabel
	ret.symTok = symTok

	return ret, true
}

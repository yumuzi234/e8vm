package asm

import (
	"e8vm.io/e8vm/arch"
	asminst "e8vm.io/e8vm/asm/inst"
	"e8vm.io/e8vm/lexing"
)

var opJmpMap = map[string]uint32{
	"j":   2,
	"jal": 3,
}

func isValidSymbol(sym string) bool {
	return true
}

func resolveInstJmp(p lexing.Logger, ops []*lexing.Token) (*inst, bool) {
	op0 := ops[0]
	opName := op0.Lit
	var op uint32

	// op sym
	switch opName {
	case "j":
		op = arch.J
	case "jal":
		op = arch.JAL
	default:
		return nil, false
	}

	var pack, sym string
	var fill int
	var symTok *lexing.Token

	if argCount(p, ops, 1) {
		symTok = ops[1]
		if checkLabel(p, ops[1]) {
			sym = ops[1].Lit
			fill = fillLabel
		} else {
			pack, sym = parseSym(p, ops[1])
			fill = fillLink
		}
	}

	ret := new(inst)
	ret.inst = asminst.Jmp(op, 0)
	ret.pkg = pack
	ret.sym = sym
	ret.fill = fill
	ret.symTok = symTok

	return ret, true
}

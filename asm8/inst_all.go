package asm8

import (
	"e8vm.io/e8vm/lex8"
)

var insts = []instResolver{
	resolveInstReg,
	resolveInstImm,
	resolveInstBr,
	resolveInstJmp,
	resolveInstSys,
}

func resolveInst(log lex8.Logger, ops []*lex8.Token) *inst {
	return instResolvers(insts).resolve(log, ops)
}

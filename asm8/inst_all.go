package asm8

import (
	"e8vm.io/e8vm/lexing"
)

var insts = []instResolver{
	resolveInstReg,
	resolveInstImm,
	resolveInstBr,
	resolveInstJmp,
	resolveInstSys,
}

func resolveInst(log lexing.Logger, ops []*lexing.Token) *inst {
	return instResolvers(insts).resolve(log, ops)
}

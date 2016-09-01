package asm

import (
	"e8vm.io/e8vm/lexing"
)

func argCount(log lexing.Logger, ops []*lexing.Token, n int) bool {
	if len(ops) == n+1 {
		return true
	}

	log.Errorf(ops[0].Pos, "%q needs %d arguments", ops[0].Lit, n)
	return false
}

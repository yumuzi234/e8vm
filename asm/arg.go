package asm

import (
	"shanhu.io/smlvm/lexing"
)

func argCount(log lexing.Logger, ops []*lexing.Token, n int) bool {
	if len(ops) == n+1 {
		return true
	}

	log.Errorf(ops[0].Pos, "%q needs %d arguments", ops[0].Lit, n)
	return false
}

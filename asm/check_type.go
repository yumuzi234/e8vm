package asm

import (
	"e8vm.io/e8vm/asm/parse"
	"e8vm.io/e8vm/lexing"
)

func checkTypeAll(p lexing.Logger, toks []*lexing.Token, typ int) bool {
	for _, t := range toks {
		if t.Type != typ {
			p.Errorf(t.Pos, "expect operand, got %s", parse.TypeStr(t.Type))
			return false
		}
	}
	return true
}

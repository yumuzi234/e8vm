package asm

import (
	"shanhu.io/smlvm/asm/parse"
	"shanhu.io/smlvm/lexing"
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

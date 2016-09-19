package asm

import (
	"shanhu.io/smlvm/lexing"
)

type inst struct {
	inst   uint32
	pkg    string
	sym    string
	fill   int
	symTok *lexing.Token
}

type instResolver func(lexing.Logger, []*lexing.Token) (*inst, bool)
type instResolvers []instResolver

func (rs instResolvers) resolve(log lexing.Logger, ops []*lexing.Token) *inst {
	for _, r := range rs {
		if i, hit := r(log, ops); hit {
			return i
		}
	}

	op0 := ops[0]
	log.Errorf(op0.Pos, "invalid asm instruction %q", op0.Lit)
	return nil
}

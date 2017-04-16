package sempass

import (
	"shanhu.io/smlvm/lexing"
)

type assigner struct {
	err      bool
	needCast bool
	mask     []bool
	pos      *lexing.Pos
}

func newAssigner(n int, pos *lexing.Pos) *assigner {
	return &assigner{
		mask: make([]bool, n),
		pos:  pos,
	}
}

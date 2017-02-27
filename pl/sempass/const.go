package sempass

import (
	"shanhu.io/smlvm/lexing"
	"shanhu.io/smlvm/pl/tast"
	"shanhu.io/smlvm/pl/types"
)

func numCast(
	b *builder, pos *lexing.Pos, v int64, from tast.Expr, to types.T,
) tast.Expr {
	if types.IsInteger(to) && types.InRange(v, to) {
		return tast.NewCast(from, to)
	}
	b.CodeErrorf(
		pos, "pl.cannotCast",
		"cannot cast const number %d to %s", v, to,
	)
	return nil
}

func numCastInt(
	b *builder, pos *lexing.Pos, v int64, from tast.Expr,
) tast.Expr {
	return numCast(b, pos, v, from, types.Int)
}

func numCastUint(
	b *builder, pos *lexing.Pos, v int64, from tast.Expr,
) tast.Expr {
	return numCast(b, pos, v, from, types.Uint)
}

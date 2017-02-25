package sempass

import (
	"shanhu.io/smlvm/lexing"
	"shanhu.io/smlvm/pl/tast"
	"shanhu.io/smlvm/pl/types"
)

func numsCast(
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

func numsCastInt(
	b *builder, pos *lexing.Pos, v int64, from tast.Expr,
) tast.Expr {
	return numsCast(b, pos, v, from, types.Int)
}

func numsCastUint(
	b *builder, pos *lexing.Pos, v int64, from tast.Expr,
) tast.Expr {
	return numsCast(b, pos, v, from, types.Uint)
}

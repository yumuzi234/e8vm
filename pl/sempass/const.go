package sempass

import (
	"shanhu.io/smlvm/lexing"
	"shanhu.io/smlvm/pl/tast"
	"shanhu.io/smlvm/pl/types"
)

func constCast(
	b *builder, pos *lexing.Pos, v int64, from tast.Expr, to types.T,
) tast.Expr {
	if types.IsInteger(to) && types.InRange(v, to) {
		return tast.NewCast(from, to)
	}
	b.Errorf(pos, "cannot cast %d to %s", v, to)
	return nil
}

func constCastInt(
	b *builder, pos *lexing.Pos, v int64, from tast.Expr,
) tast.Expr {
	return constCast(b, pos, v, from, types.Int)
}

func constCastUint(
	b *builder, pos *lexing.Pos, v int64, from tast.Expr,
) tast.Expr {
	return constCast(b, pos, v, from, types.Uint)
}

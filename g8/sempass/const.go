package sempass

import (
	"e8vm.io/e8vm/g8/tast"
	"e8vm.io/e8vm/g8/types"
	"e8vm.io/e8vm/lex8"
)

func constCast(
	b *Builder, pos *lex8.Pos, v int64, from tast.Expr, to types.T,
) tast.Expr {
	if types.IsInteger(to) && types.InRange(v, to) {
		return &tast.Cast{from, tast.NewRef(to)}
	}

	b.Errorf(pos, "cannot cast %d to %s", v, to)
	return nil
}

func constCastInt(
	b *Builder, pos *lex8.Pos, v int64, from tast.Expr,
) tast.Expr {
	return constCast(b, pos, v, from, types.Int)
}

func constCastUint(
	b *Builder, pos *lex8.Pos, v int64, from tast.Expr,
) tast.Expr {
	return constCast(b, pos, v, from, types.Uint)
}

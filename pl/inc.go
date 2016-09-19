package pl

import (
	"shanhu.io/smlvm/pl/codegen"
	"shanhu.io/smlvm/pl/tast"
)

func buildIncStmt(b *builder, stmt *tast.IncStmt) {
	expr := b.buildExpr(stmt.Expr)
	switch stmt.Op.Lit {
	case "++":
		b.b.Arith(expr.IR(), expr.IR(), "+", codegen.Num(1))
	case "--":
		b.b.Arith(expr.IR(), expr.IR(), "-", codegen.Num(1))
	default:
		panic("bug")
	}
}

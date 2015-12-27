package g8

import (
	"e8vm.io/e8vm/g8/ir"
	"e8vm.io/e8vm/g8/tast"
)

func buildIncStmt(b *builder, stmt *tast.IncStmt) {
	expr := b.buildExpr2(stmt.Expr)
	switch stmt.Op.Lit {
	case "++":
		b.b.Arith(expr.IR(), expr.IR(), "+", ir.Num(1))
	case "--":
		b.b.Arith(expr.IR(), expr.IR(), "-", ir.Num(1))
	default:
		panic("bug")
	}
}

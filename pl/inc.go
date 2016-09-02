package pl

import (
	"e8vm.io/e8vm/pl/codegen"
	"e8vm.io/e8vm/pl/tast"
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

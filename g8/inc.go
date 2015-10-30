package g8

import (
	"e8vm.io/e8vm/g8/ast"
	"e8vm.io/e8vm/g8/ir"
	"e8vm.io/e8vm/g8/types"
)

func buildIncStmt(b *builder, stmt *ast.IncStmt) {
	op := stmt.Op.Lit
	expr := buildExpr(b, stmt.Expr)
	if expr == nil {
		return
	}

	if !expr.IsSingle() {
		b.Errorf(stmt.Op.Pos, "%s on expression list", op)
		return
	}

	t := expr.Type()
	if !types.IsInteger(t) {
		b.Errorf(stmt.Op.Pos, "%s on %s", op, expr)
		return
	}

	if !expr.Addressable() {
		b.Errorf(stmt.Op.Pos, "%s on non-addressable", op)
		return
	}

	switch stmt.Op.Lit {
	case "++":
		b.b.Arith(expr.IR(), expr.IR(), "+", ir.Num(1))
	case "--":
		b.b.Arith(expr.IR(), expr.IR(), "-", ir.Num(1))
	default:
		b.Errorf(stmt.Op.Pos, "invalid inc op %s", op)
	}
}

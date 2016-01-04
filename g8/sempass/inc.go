package sempass

import (
	"e8vm.io/e8vm/g8/ast"
	"e8vm.io/e8vm/g8/tast"
	"e8vm.io/e8vm/g8/types"
)

func buildIncStmt(b *builder, stmt *ast.IncStmt) tast.Stmt {
	op := stmt.Op.Lit
	expr := b.BuildExpr(stmt.Expr)
	if expr == nil {
		return nil
	}

	ref := expr.R()
	if !ref.IsSingle() {
		b.Errorf(stmt.Op.Pos, "%s on expression list", op)
		return nil
	}

	t := ref.Type()
	if !types.IsInteger(t) {
		b.Errorf(stmt.Op.Pos, "%s on %s", op, t)
		return nil
	}

	if !ref.Addressable {
		b.Errorf(stmt.Op.Pos, "%s on non-addressable", op)
		return nil
	}

	switch stmt.Op.Lit {
	case "++", "--":
		return &tast.IncStmt{expr, stmt.Op}
	default:
		b.Errorf(stmt.Op.Pos, "invalid inc op %s", op)
		return nil
	}
}

package sempass

import (
	"shanhu.io/smlvm/pl/ast"
	"shanhu.io/smlvm/pl/tast"
	"shanhu.io/smlvm/pl/types"
)

func buildIncStmt(b *builder, stmt *ast.IncStmt) tast.Stmt {
	op := stmt.Op.Lit
	expr := b.buildExpr(stmt.Expr)
	if expr == nil {
		return nil
	}

	ref := expr.R()
	if !ref.IsSingle() {
		b.CodeErrorf(
			stmt.Op.Pos, "pl.incStmt.notSingle",
			"%s on expression list", op,
		)
		return nil
	}

	t := ref.Type()
	if !types.IsInteger(t) {
		b.CodeErrorf(
			stmt.Op.Pos, "pl.incStmt.notInt",
			"cannot %s on %s, not a integer variable", op, t,
		)
		return nil
	}

	if !ref.Addressable {
		b.CodeErrorf(
			stmt.Op.Pos, "pl.incStmt.nonAddressable",
			"%s on non-addressable", op,
		)
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

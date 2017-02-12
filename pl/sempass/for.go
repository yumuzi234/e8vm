package sempass

import (
	"shanhu.io/smlvm/pl/ast"
	"shanhu.io/smlvm/pl/tast"
)

func buildForStmt(b *builder, stmt *ast.ForStmt) tast.Stmt {
	b.scope.Push()
	defer scopePopAndCheck(b)

	ret := new(tast.ForStmt)
	if stmt.Init != nil {
		ret.Init = b.buildStmt(stmt.Init)
	}

	if stmt.Cond != nil {
		ret.Cond = b.buildExpr(stmt.Cond)
		if ret.Cond == nil {
			return nil
		}

		ref := ret.Cond.R()
		if !ref.IsBool() {
			pos := ast.ExprPos(stmt.Cond)
			b.Errorf(pos, "expect boolean expression, got %s", ref)
			return nil
		}
	}

	b.nloop++
	ret.Body = buildBlock(b, stmt.Body)
	b.nloop--

	if stmt.Iter != nil {
		ret.Iter = b.buildStmt(stmt.Iter)
	}
	return ret
}

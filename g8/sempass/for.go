package sempass

import (
	"e8vm.io/e8vm/g8/ast"
	"e8vm.io/e8vm/g8/tast"
)

func buildForStmt(b *builder, stmt *ast.ForStmt) tast.Stmt {
	b.scope.Push()
	defer b.scope.Pop()

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

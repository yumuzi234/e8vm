package sempass

import (
	"e8vm.io/e8vm/pl/ast"
	"e8vm.io/e8vm/pl/tast"
)

func buildContinueStmt(b *builder, s *ast.ContinueStmt) tast.Stmt {
	if s.Label != nil {
		b.Errorf(s.Label.Pos, "continue with label not implemented")
		return nil
	}
	if b.nloop == 0 {
		b.Errorf(s.Kw.Pos, "continue is not in a for block")
		return nil
	}
	return &tast.ContinueStmt{}
}

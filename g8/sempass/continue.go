package sempass

import (
	"e8vm.io/e8vm/g8/ast"
	"e8vm.io/e8vm/g8/tast"
)

func buildContinueStmt(b *Builder, s *ast.ContinueStmt) *tast.ContinueStmt {
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

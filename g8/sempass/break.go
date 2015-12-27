package sempass

import (
	"e8vm.io/e8vm/g8/ast"
	"e8vm.io/e8vm/g8/tast"
)

func buildBreakStmt(b *Builder, s *ast.BreakStmt) tast.Stmt {
	if s.Label != nil {
		b.Errorf(s.Label.Pos, "break with label not implemented")
		return nil
	}
	if b.nloop == 0 {
		b.Errorf(s.Kw.Pos, "break is not in a for block")
		return nil
	}
	return &tast.BreakStmt{}
}

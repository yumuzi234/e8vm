package sempass

import (
	"shanhu.io/smlvm/pl/ast"
	"shanhu.io/smlvm/pl/tast"
)

func buildContinueStmt(b *builder, s *ast.ContinueStmt) tast.Stmt {
	if s.Label != nil {
		b.Errorf(s.Label.Pos, "continue with label not implemented")
		return nil
	}
	if b.nloop == 0 {
		b.CodeErrorf(s.Kw.Pos, "pl.continueStmt.notInLoop",
			"continue is not in a for block")
		return nil
	}
	return &tast.ContinueStmt{}
}

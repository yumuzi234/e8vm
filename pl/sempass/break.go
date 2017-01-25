package sempass

import (
	"shanhu.io/smlvm/pl/ast"
	"shanhu.io/smlvm/pl/tast"
)

func buildBreakStmt(b *builder, s *ast.BreakStmt) tast.Stmt {
	if s.Label != nil {
		b.CodeErrorf(s.Label.Pos, "pl.notYetSupported",
			"break with label not implemented")
		return nil
	}
	if b.nloop == 0 {
		b.CodeErrorf(s.Kw.Pos, "pl.breakStmt.notInLoop",
			"break is not in a for block")
		return nil
	}
	return &tast.BreakStmt{}
}

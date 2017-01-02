package sempass

import (
	"shanhu.io/smlvm/pl/ast"
	"shanhu.io/smlvm/pl/tast"
)

func scopePopAndCheck(b *builder) {
	tab := b.scope.Pop()
	syms := tab.List()
	//Same error code used here, instead of different codes based on sym.Type
	for _, sym := range syms {
		if !sym.Used {
			b.CodeErrorf(
				sym.Pos, "pl.unusedSym",
				"unused %s %q", tast.SymStr(sym.Type), sym.Name(),
			)
		}
	}
}

func buildBlock(b *builder, block *ast.Block) tast.Stmt {
	b.scope.Push()
	defer scopePopAndCheck(b)

	var stmts []tast.Stmt
	for _, stmt := range block.Stmts {
		s := b.buildStmt(stmt)
		if s != nil {
			stmts = append(stmts, s)
		}
	}

	return &tast.Block{stmts}
}

func buildBlockStmt(b *builder, block *ast.BlockStmt) tast.Stmt {
	return buildBlock(b, block.Block)
}

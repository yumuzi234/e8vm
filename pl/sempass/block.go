package sempass

import (
	"e8vm.io/e8vm/pl/ast"
	"e8vm.io/e8vm/pl/tast"
)

func scopePopAndCheck(b *builder) {
	tab := b.scope.Pop()
	syms := tab.List()
	for _, sym := range syms {
		if !sym.Used {
			b.Errorf(
				sym.Pos,
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

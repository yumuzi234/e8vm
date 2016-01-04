package sempass

import (
	"e8vm.io/e8vm/g8/ast"
	"e8vm.io/e8vm/g8/tast"
)

func buildBlock(b *builder, block *ast.Block) tast.Stmt {
	b.scope.Push()
	defer b.scope.Pop()

	var stmts []tast.Stmt
	for _, stmt := range block.Stmts {
		s := b.BuildStmt(stmt)
		if s != nil {
			stmts = append(stmts, s)
		}
	}

	return &tast.Block{stmts}
}

func buildBlockStmt(b *builder, block *ast.BlockStmt) tast.Stmt {
	return buildBlock(b, block.Block)
}

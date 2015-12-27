package g8

import (
	"e8vm.io/e8vm/g8/ast"
	"e8vm.io/e8vm/g8/tast"
)

func buildBlock(b *builder, stmt *ast.Block) {
	b.scope.Push()
	defer b.scope.Pop()

	b.buildStmts(stmt.Stmts)
}

func genBlock(b *builder, block *tast.Block) {
	for _, s := range block.Stmts {
		b.buildStmt2(s)
	}
}

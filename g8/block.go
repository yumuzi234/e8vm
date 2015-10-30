package g8

import (
	"e8vm.io/e8vm/g8/ast"
)

func buildBlock(b *builder, stmt *ast.Block) {
	b.scope.Push()
	defer b.scope.Pop()

	b.buildStmts(stmt.Stmts)
}

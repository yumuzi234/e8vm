package g8

import (
	"e8vm.io/e8vm/g8/tast"
)

func genBlock(b *builder, block *tast.Block) {
	for _, s := range block.Stmts {
		b.buildStmt2(s)
	}
}

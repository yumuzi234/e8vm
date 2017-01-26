package pl

import (
	"fmt"

	"shanhu.io/smlvm/pl/tast"
)

func buildSwitchStmt(b *builder, stmt *tast.SwitchStmt) {
	fmt.Println("switch")
}

func buildFallthroughStmt(b *builder) {
}

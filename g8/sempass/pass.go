package sempass

import (
	"e8vm.io/e8vm/sym8"
)

// NewBuilder creates a new builder with a specific path.
func NewBuilder(path string, scope *sym8.Scope) *Builder {
	ret := newBuilder(path)
	ret.exprFunc = buildExpr
	ret.constFunc = buildConstExpr
	ret.typeFunc = buildType
	ret.stmtFunc = buildStmt

	ret.scope = scope // TODO: remove this

	return ret
}

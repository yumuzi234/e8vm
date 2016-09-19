package sempass

import (
	"shanhu.io/smlvm/lexing"
	"shanhu.io/smlvm/pl/ast"
	"shanhu.io/smlvm/pl/tast"
	"shanhu.io/smlvm/syms"
)

// BuildBareFunc build a list of statements.
func BuildBareFunc(scope *syms.Scope, stmts []ast.Stmt) (
	[]tast.Stmt, []*lexing.Error,
) {
	b := makeBuilder("_", scope)
	b.scope.Push()
	defer b.scope.Pop()
	ret := buildStmts(b, stmts)
	errs := b.Errs()
	if errs != nil {
		return nil, errs
	}
	return ret, nil
}

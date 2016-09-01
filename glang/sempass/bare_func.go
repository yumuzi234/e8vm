package sempass

import (
	"e8vm.io/e8vm/glang/ast"
	"e8vm.io/e8vm/glang/tast"
	"e8vm.io/e8vm/lexing"
	"e8vm.io/e8vm/syms"
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

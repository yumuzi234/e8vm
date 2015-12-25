package sempass

import (
	"e8vm.io/e8vm/dagvis"
	"e8vm.io/e8vm/g8/ast"
	"e8vm.io/e8vm/g8/tast"
	"e8vm.io/e8vm/lex8"
	"e8vm.io/e8vm/sym8"
)

// Builder is a semantic pass builder.
type Builder struct {
	*lex8.ErrorList
	path string

	this  *tast.Ref
	scope *sym8.Scope

	exprFunc  func(b *Builder, expr ast.Expr) tast.Expr
	constFunc func(b *Builder, expr ast.Expr) tast.Expr
	stmtFunc  func(b *Builder, stmt ast.Stmt) tast.Stmt

	// file level dependency, for checking circular dependencies.
	deps deps
}

func newBuilder(path string) *Builder {
	return &Builder{
		ErrorList: lex8.NewErrorList(),
		path:      path,
		scope:     sym8.NewScope(),
	}
}

// BuildExpr builds the expression.
func (b *Builder) BuildExpr(expr ast.Expr) tast.Expr {
	return b.exprFunc(b, expr)
}

// BuildConstExpr builds a constant expression.
func (b *Builder) BuildConstExpr(expr ast.Expr) tast.Expr {
	return b.constFunc(b, expr)
}

// RefSym references a symbol.
func (b *Builder) RefSym(sym *sym8.Symbol, pos *lex8.Pos) {
	// track file dependencies inside a package
	if b.deps == nil {
		return // no need to track deps
	}

	symPos := sym.Pos
	if symPos == nil {
		return // builtin
	}
	if sym.Pkg() != b.path {
		return // cross package reference
	}
	if pos.File == symPos.File {
		return
	}

	b.deps.add(pos.File, symPos.File)
}

// InitDeps initializes the dependency graph.
func (b *Builder) InitDeps(asts map[string]*ast.File) {
	b.deps = newDeps(asts)
}

// SetThis sets the reference for this keyword.
func (b *Builder) SetThis(ref *tast.Ref) {
	b.this = ref
}

// DepGraph returns the dependency graph.
func (b *Builder) DepGraph() *dagvis.Graph { return b.deps.graph() }

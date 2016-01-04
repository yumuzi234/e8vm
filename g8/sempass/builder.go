package sempass

import (
	"e8vm.io/e8vm/dagvis"
	"e8vm.io/e8vm/g8/ast"
	"e8vm.io/e8vm/g8/tast"
	"e8vm.io/e8vm/g8/types"
	"e8vm.io/e8vm/lex8"
	"e8vm.io/e8vm/sym8"
)

// Builder is a semantic pass builder.
type builder struct {
	*lex8.ErrorList
	path string

	scope *sym8.Scope

	exprFunc  func(b *builder, expr ast.Expr) tast.Expr
	constFunc func(b *builder, expr ast.Expr) tast.Expr
	stmtFunc  func(b *builder, stmt ast.Stmt) tast.Stmt
	typeFunc  func(b *builder, expr ast.Expr) types.T

	// file level dependency, for checking circular dependencies.
	deps deps

	nloop    int
	this     *tast.Ref
	thisType *types.Pointer

	retType  []types.T
	retNamed bool
}

func newBuilder(path string) *builder {
	return &builder{
		ErrorList: lex8.NewErrorList(),
		path:      path,
		scope:     sym8.NewScope(),
	}
}

// BuildExpr builds the expression.
func (b *builder) BuildExpr(expr ast.Expr) tast.Expr {
	return b.exprFunc(b, expr)
}

func (b *builder) buildConstExpr(expr ast.Expr) tast.Expr {
	return b.constFunc(b, expr)
}

// BuildConstExpr builds a constant expression.
func (b *builder) BuildConstExpr(expr ast.Expr) *tast.Const {
	ret, ok := b.buildConstExpr(expr).(*tast.Const)
	if !ok {
		b.Errorf(ast.ExprPos(expr), "expect a const")
		return nil
	}
	return ret
}

// BuildType builds an expression that represents a type.
func (b *builder) BuildType(expr ast.Expr) types.T {
	return b.typeFunc(b, expr)
}

// BuildStmt builds a statement.
func (b *builder) BuildStmt(stmt ast.Stmt) tast.Stmt {
	return b.stmtFunc(b, stmt)
}

// RefSym references a symbol.
func (b *builder) RefSym(sym *sym8.Symbol, pos *lex8.Pos) {
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
func (b *builder) InitDeps(asts map[string]*ast.File) {
	b.deps = newDeps(asts)
}

// SetThis sets the reference for this keyword.
func (b *builder) SetThis(ref *tast.Ref) { b.this = ref }

// SetRetType sets the return type of the current function.
func (b *builder) SetRetType(ts []types.T, named bool) {
	b.retType = ts
	b.retNamed = named
}

// DepGraph returns the dependency graph.
func (b *builder) DepGraph() *dagvis.Graph { return b.deps.graph() }

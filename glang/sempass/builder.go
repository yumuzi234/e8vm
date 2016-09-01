package sempass

import (
	"e8vm.io/e8vm/dagvis"
	"e8vm.io/e8vm/glang/ast"
	"e8vm.io/e8vm/glang/tast"
	"e8vm.io/e8vm/glang/types"
	"e8vm.io/e8vm/lexing"
	"e8vm.io/e8vm/sym8"
)

type builder struct {
	*lexing.ErrorList
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

	// if the parsing is in left hand side.
	// when in left hand side, referecing a variable does not count.
	lhs bool
}

func newBuilder(path string, scope *sym8.Scope) *builder {
	return &builder{
		ErrorList: lexing.NewErrorList(),
		path:      path,
		scope:     scope,
	}
}

func (b *builder) lhsSwap(lhs bool) bool {
	lhs, b.lhs = b.lhs, lhs
	return lhs
}

func (b *builder) lhsRestore(lhs bool) { b.lhs = lhs }

func (b *builder) buildExpr(expr ast.Expr) tast.Expr {
	return b.exprFunc(b, expr)
}

func (b *builder) buildConstExpr(expr ast.Expr) tast.Expr {
	return b.constFunc(b, expr)
}

func (b *builder) buildConst(expr ast.Expr) *tast.Const {
	ret, ok := b.buildConstExpr(expr).(*tast.Const)
	if !ok {
		b.Errorf(ast.ExprPos(expr), "expect a const")
		return nil
	}
	return ret
}

func (b *builder) buildType(expr ast.Expr) types.T {
	return b.typeFunc(b, expr)
}

func (b *builder) buildStmt(stmt ast.Stmt) tast.Stmt {
	return b.stmtFunc(b, stmt)
}

func (b *builder) refSym(sym *sym8.Symbol, pos *lexing.Pos) {
	if !b.lhs {
		sym.Used = true
	}

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

func (b *builder) initDeps(asts map[string]*ast.File) {
	b.deps = newDeps(asts)
}

func (b *builder) depGraph() *dagvis.Graph { return b.deps.graph() }

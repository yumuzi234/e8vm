package g8

import (
	"fmt"
	"io"

	"e8vm.io/e8vm/g8/ir"
	"e8vm.io/e8vm/g8/tast"
	"e8vm.io/e8vm/g8/types"
	"e8vm.io/e8vm/lex8"
	"e8vm.io/e8vm/sym8"
)

// builder builds a package
type builder struct {
	*lex8.ErrorList
	path  string
	scope *sym8.Scope

	p *ir.Pkg
	f *ir.Func
	b *ir.Block

	panicFunc ir.Ref // for calling panic
	fretRef   *ref   // to store return value
	this      *ref   // not nil when building a method

	continues *blockStack
	breaks    *blockStack

	exprFunc func(b *builder, expr tast.Expr) *ref
	stmtFunc func(b *builder, stmt tast.Stmt)
	irLog    io.WriteCloser

	anonyCount int
}

func newBuilder(path string) *builder {
	s := sym8.NewScope()
	return &builder{
		ErrorList: lex8.NewErrorList(),
		path:      path,
		p:         ir.NewPkg(path),
		scope:     s, // package scope

		continues: newBlockStack(),
		breaks:    newBlockStack(),
	}
}

func (b *builder) anonyName(name string) string {
	if name == "_" {
		name = fmt.Sprintf("_:%d", b.anonyCount)
		b.anonyCount++
	}
	return name
}

func (b *builder) newTempIR(t types.T) ir.Ref {
	return b.f.NewTemp(t.Size(), types.IsByte(t), t.RegSizeAlign())
}

func (b *builder) newTemp(t types.T) *ref { return newRef(t, b.newTempIR(t)) }

func (b *builder) newCond() ir.Ref { return b.f.NewTemp(1, true, false) }
func (b *builder) newPtr() ir.Ref  { return b.f.NewTemp(4, true, true) }

func (b *builder) newAddressableTemp(t types.T) *ref {
	return newAddressableRef(t, b.newTempIR(t))
}

func (b *builder) newLocal(t types.T, name string) ir.Ref {
	return b.f.NewLocal(t.Size(), name,
		types.IsByte(t), t.RegSizeAlign(),
	)
}

func (b *builder) newGlobalVar(t types.T, name string) ir.Ref {
	name = b.anonyName(name)
	return b.p.NewGlobalVar(t.Size(), name, types.IsByte(t), t.RegSizeAlign())
}

func (b *builder) buildExpr(expr tast.Expr) *ref {
	if b.exprFunc == nil {
		panic("exprFunc2 function not assigned")
	}
	return b.exprFunc(b, expr)
}

func (b *builder) buildStmt(stmt tast.Stmt) { b.stmtFunc(b, stmt) }

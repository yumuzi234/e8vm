package g8

import (
	"fmt"

	"e8vm.io/e8vm/g8/codegen"
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

	p *codegen.Pkg
	f *codegen.Func
	b *codegen.Block

	panicFunc codegen.Ref // for calling panic
	fretRef   *ref        // to store return value
	this      *ref        // not nil when building a method

	continues *blockStack
	breaks    *blockStack

	exprFunc func(b *builder, expr tast.Expr) *ref
	stmtFunc func(b *builder, stmt tast.Stmt)

	anonyCount int // count for "_"
}

func newBuilder(path string) *builder {
	s := sym8.NewScope()
	return &builder{
		ErrorList: lex8.NewErrorList(),
		path:      path,
		p:         codegen.NewPkg(path),
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

func (b *builder) newTempIR(t types.T) codegen.Ref {
	return b.f.NewTemp(t.Size(), types.IsByte(t), t.RegSizeAlign())
}

func (b *builder) newTemp(t types.T) *ref { return newRef(t, b.newTempIR(t)) }

func (b *builder) newCond() codegen.Ref { return b.f.NewTemp(1, true, false) }
func (b *builder) newPtr() codegen.Ref  { return b.f.NewTemp(4, true, true) }

func (b *builder) newAddressableTemp(t types.T) *ref {
	return newAddressableRef(t, b.newTempIR(t))
}

func (b *builder) newLocal(t types.T, name string) codegen.Ref {
	return b.f.NewLocal(t.Size(), name,
		types.IsByte(t), t.RegSizeAlign(),
	)
}

func (b *builder) newGlobalVar(t types.T, name string) codegen.Ref {
	name = b.anonyName(name)
	return b.p.NewGlobalVar(t.Size(), name, types.IsByte(t), t.RegSizeAlign())
}

func (b *builder) buildExpr(expr tast.Expr) *ref {
	return b.exprFunc(b, expr)
}

func (b *builder) buildStmt(stmt tast.Stmt) { b.stmtFunc(b, stmt) }

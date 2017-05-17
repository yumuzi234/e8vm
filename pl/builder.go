package pl

import (
	"fmt"

	"shanhu.io/smlvm/lexing"
	"shanhu.io/smlvm/pl/codegen"
	"shanhu.io/smlvm/pl/tast"
	"shanhu.io/smlvm/pl/types"
	"shanhu.io/smlvm/syms"
)

// builder builds a package
type builder struct {
	*lexing.ErrorList
	path  string
	scope *syms.Scope

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

	vTableMap map[*types.Interface]*vTable
}

func (b *builder) newImplement(i *types.Interface, r *ref, s *types.Struct) {
	t := b.vTableMap[i]
	if t == nil {
		t = newTable(i)
		b.vTableMap[i] = t
	}
	if t.implementMap[s] == nil {
		slice := make([]*syms.Symbol, len(t.funcs))
		for i, funcName := range t.funcs {
			// check how to change []*Symbol to []*ref
			slice[i] = s.Syms.Query(funcName)
		}
		t.implementMap[s] = slice
	}
}

func newBuilder(path string) *builder {
	s := syms.NewScope()
	return &builder{
		ErrorList: lexing.NewErrorList(),
		path:      path,
		p:         codegen.NewPkg(path),
		scope:     s, // package scope

		continues: newBlockStack(),
		breaks:    newBlockStack(),

		vTableMap: make(map[*types.Interface]*vTable),
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

package glang

import (
	"e8vm.io/e8vm/asm8"
	"e8vm.io/e8vm/glang/codegen"
	"e8vm.io/e8vm/glang/tast"
	"e8vm.io/e8vm/glang/types"
)

func buildPkgSym(b *builder, m *tast.MemberExpr, pkg *types.Pkg) *ref {
	sym := m.Symbol
	if pkg.Lang == "asm8" {
		switch sym.Type {
		case asm8.SymVar:
			if m.Type() != types.Uint {
				panic("bug")
			}
			ptr := b.newTemp(types.Uint)
			s := codegen.NewHeapSym(sym.Pkg(), sym.Name(), 0, false, false)
			b.b.Arith(ptr.IR(), nil, "&", s)
			return ptr
		case asm8.SymFunc:
			if !types.SameType(m.Type(), types.VoidFunc) {
				panic("bug")
			}
			return newRef(
				types.VoidFunc,
				codegen.NewFuncSym(
					sym.Pkg(), sym.Name(), codegen.VoidFuncSig,
				),
			)
		}
		panic("bug")
	}

	switch sym.Type {
	case tast.SymConst:
		return sym.Obj.(*objConst).ref
	case tast.SymVar:
		return sym.Obj.(*objVar).ref
	case tast.SymFunc:
		return sym.Obj.(*objFunc).ref
	}
	panic("bug")
}

func buildMember(b *builder, m *tast.MemberExpr) *ref {
	obj := b.buildExpr(m.Expr)
	if !obj.IsSingle() {
		panic("not single")
	}

	t := obj.Type()
	if pkg, ok := t.(*types.Pkg); ok {
		return buildPkgSym(b, m, pkg)
	}

	pt := types.PointerOf(t)
	var tstruct *types.Struct
	if pt != nil {
		tstruct = pt.(*types.Struct)
	} else {
		tstruct = t.(*types.Struct)
	}

	addr := b.newPtr()
	if pt != nil {
		b.b.Assign(addr, obj.IR())
		if obj != b.this {
			nilPointerPanic(b, addr)
		}
	} else {
		b.b.Arith(addr, nil, "&", obj.IR())
	}

	sym := m.Symbol
	if sym.Type == tast.SymField {
		return buildField(b, addr, sym.Obj.(*types.Field))
	} else if sym.Type == tast.SymFunc {
		recv := newRef(types.NewPointer(tstruct), addr)
		method := sym.Obj.(*objFunc)
		ft := method.Type().(*types.Func)
		return newRecvRef(ft, recv, method.IR())
	}

	panic("bug")
}

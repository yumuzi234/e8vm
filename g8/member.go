package g8

import (
	"e8vm.io/e8vm/asm8"
	"e8vm.io/e8vm/g8/ast"
	"e8vm.io/e8vm/g8/ir"
	"e8vm.io/e8vm/g8/types"
	"e8vm.io/e8vm/lex8"
	"e8vm.io/e8vm/sym8"
)

func findPackageSym(
	b *builder, sub *lex8.Token, pkg *types.Pkg,
) *sym8.Symbol {
	sym := pkg.Syms.Query(sub.Lit)
	if sym == nil {
		b.Errorf(sub.Pos, "%s has no symbol named %s",
			pkg, sub.Lit,
		)
		return nil
	}
	name := sym.Name()
	if !sym8.IsPublic(name) && sym.Pkg() != b.path {
		b.Errorf(sub.Pos, "symbol %s is not public", name)
		return nil
	}

	return sym
}

func buildPackageSym(b *builder, m *ast.MemberExpr, pkg *types.Pkg) *ref {
	sym := findPackageSym(b, m.Sub, pkg)
	if sym == nil {
		return nil
	}

	if pkg.Lang == "asm8" {
		switch sym.Type {
		case asm8.SymVar:
			ptr := b.newTemp(types.Uint)
			s := ir.NewHeapSym(sym.Pkg(), sym.Name(), 0, false, false)
			b.b.Arith(ptr.IR(), nil, "&", s)
			return ptr
		case asm8.SymFunc:
			return newRef(
				types.VoidFunc,
				ir.NewFuncSym(sym.Pkg(), sym.Name(), ir.VoidFuncSig),
			)
		}

		b.Errorf(m.Sub.Pos, "invalid symbol %s in %s: %s",
			m.Sub.Lit, pkg, asm8.SymStr(sym.Type),
		)
		return nil
	}

	switch sym.Type {
	case symConst:
		return sym.Obj.(*objConst).ref
	case symVar:
		return sym.Obj.(*objVar).ref
	case symStruct:
		return sym.Obj.(*objType).ref
	case symFunc:
		return sym.Obj.(*objFunc).ref
	}

	b.Errorf(m.Sub.Pos, "bug: invalid symbol %s in %s: %s",
		m.Sub.Lit, pkg, symStr(sym.Type),
	)
	return nil
}

func buildConstMember(b *builder, m *ast.MemberExpr) *ref {
	obj := b.buildConstExpr(m.Expr)
	if obj == nil {
		return nil
	}
	if !obj.IsSingle() {
		b.Errorf(m.Dot.Pos, "expression list does not have any member")
		return nil
	}

	if pkg, ok := obj.Type().(*types.Pkg); ok {
		s := findPackageSym(b, m.Sub, pkg)
		if s == nil {
			return nil
		}
		switch s.Type {
		case symConst:
			return s.Obj.(*objConst).ref
		case symStruct:
			return s.Obj.(*objType).ref
		}

		b.Errorf(m.Sub.Pos, "%s.%s is not a const", pkg, m.Sub.Lit)
		return nil
	}

	b.Errorf(m.Dot.Pos, "expect const expression")
	return nil
}

func buildMember(b *builder, m *ast.MemberExpr) *ref {
	obj := b.buildExpr(m.Expr)
	if obj == nil {
		return nil
	}

	if !obj.IsSingle() {
		b.Errorf(m.Dot.Pos, "expression list does not have any member")
		return nil
	}

	if pkg, ok := obj.Type().(*types.Pkg); ok {
		return buildPackageSym(b, m, pkg)
	}

	t := obj.Type()
	pt := types.PointerOf(t)

	var tstruct *types.Struct
	var ok bool

	if pt != nil {
		tstruct, ok = pt.(*types.Struct)
		if !ok {
			b.Errorf(m.Dot.Pos, "*%s is not a pointer of struct", t)
			return nil
		}
	} else {
		tstruct, ok = t.(*types.Struct)
		if !ok {
			b.Errorf(m.Dot.Pos, "%s is not a struct", t)
			return nil
		}
	}

	addr := b.newPtr()
	if pt != nil {
		b.b.Assign(addr, obj.IR())
		if obj != b.this {
			nilPointerPanic(b, addr)
		}
	} else {
		b.b.Arith(addr, nil, "&", obj.IR()) // load address
	}

	// TODO: the use of structFields here is more like a hack. It should not
	// be the way to judge if a struct belongs to another package or not.
	symTable := b.structFields[tstruct]
	if symTable == nil {
		// in this package.
		symTable = tstruct.Syms
	}
	name := m.Sub.Lit
	sym := symTable.Query(name)
	if sym == nil {
		b.Errorf(m.Sub.Pos, "struct %s has no member named %s",
			tstruct, m.Sub.Lit,
		)
		return nil
	} else if !sym8.IsPublic(name) && sym.Pkg() != b.path {
		b.Errorf(m.Sub.Pos, "symbol %s is not public", name)
		return nil
	}

	b.refSym(sym, m.Sub.Pos)

	if sym.Type == symField {
		return buildField(b, addr, sym.Obj.(*objField).Field)
	} else if sym.Type == symFunc {
		recv := newRef(types.NewPointer(tstruct), addr)
		method := sym.Obj.(*objFunc)
		ft := method.Type().(*types.Func)
		return newRecvRef(ft, recv, method.IR())
	}

	b.Errorf(m.Sub.Pos, "invalid sym type", symStr(sym.Type))
	return nil
}

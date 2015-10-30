package g8

import (
	"e8vm.io/e8vm/g8/ast"
	"e8vm.io/e8vm/g8/types"
	"e8vm.io/e8vm/sym8"
)

func buildPackageSym(b *builder, m *ast.MemberExpr, pkg *types.Pkg) *ref {
	sym := pkg.Syms.Query(m.Sub.Lit)
	if sym == nil {
		b.Errorf(m.Sub.Pos, "%s has no symbol named %s",
			pkg, m.Sub.Lit,
		)
		return nil
	}
	name := sym.Name()
	if !sym8.IsPublic(name) && sym.Pkg() != b.symPkg {
		b.Errorf(m.Sub.Pos, "symbol %s is not public", name)
		return nil
	}

	switch sym.Type {
	case symConst:
		return sym.Item.(*objConst).ref
	case symVar:
		return sym.Item.(*objVar).ref
	case symStruct:
		return sym.Item.(*objType).ref
	case symFunc:
		return sym.Item.(*objFunc).ref
	}

	b.Errorf(m.Sub.Pos, "bug: invalid symbol %s in %s: %s",
		m.Sub.Lit, pkg, symStr(sym.Type),
	)
	return nil
}

func buildConstMember(b *builder, m *ast.MemberExpr) *ref {
	b.Errorf(m.Dot.Pos, "const member not implemented")
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

	symTable := b.structFields[tstruct]
	if symTable == nil {
		symTable = tstruct.Syms
	}
	name := m.Sub.Lit
	sym := symTable.Query(name)
	if sym == nil {
		b.Errorf(m.Sub.Pos, "struct %s has no member named %s",
			tstruct, m.Sub.Lit,
		)
		return nil
	} else if !sym8.IsPublic(name) && sym.Pkg() != b.symPkg {
		b.Errorf(m.Sub.Pos, "symbol %s is not public", name)
		return nil
	}

	b.refSym(sym, m.Sub.Pos)

	if sym.Type == symField {
		return buildField(b, addr, sym.Item.(*objField).Field)
	} else if sym.Type == symFunc {
		recv := newRef(types.NewPointer(tstruct), addr)
		method := sym.Item.(*objFunc)
		ft := method.Type().(*types.Func)
		return newRecvRef(ft, recv, method.IR())
	}

	b.Errorf(m.Sub.Pos, "invalid sym type", symStr(sym.Type))
	return nil
}

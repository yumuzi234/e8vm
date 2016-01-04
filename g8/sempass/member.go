package sempass

import (
	"e8vm.io/e8vm/asm8"
	"e8vm.io/e8vm/g8/ast"
	"e8vm.io/e8vm/g8/tast"
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

func buildConstMember(b *builder, m *ast.MemberExpr) tast.Expr {
	obj := b.buildConstExpr(m.Expr)
	if obj == nil {
		return nil
	}

	ref := obj.R()
	if !ref.IsSingle() {
		b.Errorf(m.Dot.Pos, "%s does not have any member", ref)
		return nil
	}

	if pkg, ok := ref.T.(*types.Pkg); ok {
		s := findPackageSym(b, m.Sub, pkg)
		if s == nil {
			return nil
		}
		if s.Type != tast.SymConst {
			b.Errorf(m.Sub.Pos, "%s.%s is not a const", pkg, m.Sub.Lit)
			return nil
		}
		return &tast.Const{tast.NewRef(s.ObjType.(types.T))}
	}

	b.Errorf(m.Dot.Pos, "expect const expression")
	return nil
}

func buildPkgSym(
	b *builder, m *ast.MemberExpr, pkg *types.Pkg,
) (*tast.Ref, *sym8.Symbol) {
	sym := findPackageSym(b, m.Sub, pkg)
	if sym == nil {
		return nil, nil
	}

	if pkg.Lang == "asm8" {
		switch sym.Type {
		case asm8.SymVar:
			return tast.NewRef(types.Uint), sym
		case asm8.SymFunc:
			return tast.NewRef(types.VoidFunc), sym
		}

		b.Errorf(m.Sub.Pos, "invalid symbol %s in %s: %s",
			m.Sub.Lit, pkg, asm8.SymStr(sym.Type),
		)
		return nil, nil
	}
	t := sym.ObjType.(types.T)
	switch sym.Type {
	case tast.SymConst, tast.SymStruct, tast.SymFunc:
		return tast.NewRef(t), sym
	case tast.SymVar:
		return tast.NewAddressableRef(t), sym
	}

	b.Errorf(m.Sub.Pos, "bug: invalid symbol %s in %s: %s",
		m.Sub.Lit, pkg, tast.SymStr(sym.Type),
	)
	return nil, nil
}

func buildMember(b *builder, m *ast.MemberExpr) tast.Expr {
	obj := b.buildExpr(m.Expr)
	if obj == nil {
		return nil
	}

	ref := obj.R()
	if !ref.IsSingle() {
		b.Errorf(m.Dot.Pos, "%s does not have any member", ref)
		return nil
	}

	t := ref.T
	if pkg, ok := t.(*types.Pkg); ok {
		r, sym := buildPkgSym(b, m, pkg)
		if r == nil {
			return nil
		}
		// TODO: this can be further optimized
		return &tast.MemberExpr{obj, m.Sub, r, sym}
	}

	pt := types.PointerOf(t)
	var tstruct *types.Struct
	var ok bool
	if pt != nil {
		if tstruct, ok = pt.(*types.Struct); !ok {
			b.Errorf(m.Dot.Pos, "*%s is not a pointer of struct", t)
			return nil
		}
	} else {
		if tstruct, ok = t.(*types.Struct); !ok {
			b.Errorf(m.Dot.Pos, "%s is not a struct", t)
			return nil
		}
	}

	symTable := tstruct.Syms
	name := m.Sub.Lit
	sym := symTable.Query(name)
	if sym == nil {
		b.Errorf(m.Sub.Pos, "struct %s has no member named %s",
			tstruct, name,
		)
		return nil
	} else if !sym8.IsPublic(name) && sym.Pkg() != b.path {
		b.Errorf(m.Sub.Pos, "symbol %s is not public", name)
		return nil
	}

	b.refSym(sym, m.Sub.Pos)

	if sym.Type == tast.SymField {
		t := sym.ObjType.(types.T)
		r := tast.NewAddressableRef(t)
		return &tast.MemberExpr{obj, m.Sub, r, sym}
	} else if sym.Type == tast.SymFunc {
		ft := sym.ObjType.(*types.Func)
		r := tast.NewRef(ft.MethodFunc)
		r.Recv = ref
		return &tast.MemberExpr{obj, m.Sub, r, sym}
	}

	b.Errorf(m.Sub.Pos, "invalid sym type: %s", tast.SymStr(sym.Type))
	return nil
}

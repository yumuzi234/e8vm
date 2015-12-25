package g8

import (
	"e8vm.io/e8vm/g8/ast"
	"e8vm.io/e8vm/g8/tast"
	"e8vm.io/e8vm/g8/types"
	"e8vm.io/e8vm/lex8"
	"e8vm.io/e8vm/sym8"
)

func buildArrayType(b *builder, expr *ast.ArrayTypeExpr) types.T {
	t := buildType(b, expr.Type)
	if t == nil {
		return nil
	}

	if expr.Len == nil {
		// slice
		return &types.Slice{t}
	}

	// array
	n := b.buildConstExpr(expr.Len)
	if n == nil {
		return nil
	} else if !n.IsSingle() {
		panic("bug")
	}

	ntype := n.Type()
	c, ok := ntype.(*types.Const)
	if !ok {
		// might be true, false, or other builtin consts
		b.Errorf(ast.ExprPos(expr), "array index is not a constant")
		return nil
	}

	if v, ok := types.NumConst(ntype); ok {
		if v < 0 {
			b.Errorf(ast.ExprPos(expr),
				"array index is negative: %d", c.Value,
			)
			return nil
		} else if !types.InRange(v, types.Int) {
			b.Errorf(ast.ExprPos(expr), "index out of range of int32")
			return nil
		}
		return &types.Array{T: t, N: int32(v)}
	}

	// TODO: support typed const
	b.Errorf(ast.ExprPos(expr), "typed const not implemented yet")
	return nil
}

func buildPkgRef(b *builder, ident *lex8.Token) *types.Pkg {
	s := b.scope.Query(ident.Lit)
	if s == nil {
		b.Errorf(ident.Pos, "undefined identifier %s", ident.Lit)
		return nil
	}

	b.spass.RefSym(s, ident.Pos)
	if s.Type != tast.SymImport {
		b.Errorf(ident.Pos, "%s is not an imported package", ident.Lit)
		return nil
	}

	return s.Obj.(*objImport).ref.Type().(*types.Pkg)
}

func buildType(b *builder, expr ast.Expr) types.T {
	if expr == nil {
		panic("bug")
	}

	switch expr := expr.(type) {
	case *ast.Operand:
		ret := b.buildExpr(expr)
		if ret == nil {
			return nil
		} else if !ret.IsType() {
			b.Errorf(ast.ExprPos(expr), "expect a type, got %s", ret)
			return nil
		}
		return ret.TypeType()
	case *ast.StarExpr:
		t := buildType(b, expr.Expr)
		if t == nil {
			return nil
		}
		return &types.Pointer{t}
	case *ast.ArrayTypeExpr:
		return buildArrayType(b, expr)
	case *ast.ParenExpr:
		return buildType(b, expr.Expr)
	case *ast.FuncTypeExpr:
		return buildFuncType(b, nil, expr.FuncSig)
	case *ast.MemberExpr:
		op, ok := expr.Expr.(*ast.Operand)
		if !ok {
			b.Errorf(ast.ExprPos(expr.Expr), "expect a package")
			return nil
		}
		pkg := buildPkgRef(b, op.Token)
		if pkg == nil {
			return nil
		}
		name := expr.Sub.Lit
		s := pkg.Syms.Query(name)
		if s == nil {
			b.Errorf(expr.Sub.Pos, "symbol %s not found", name)
			return nil
		}
		if !sym8.IsPublic(name) && s.Pkg() != b.path {
			b.Errorf(expr.Sub.Pos, "symbol %s is not public", name)
			return nil
		}

		if s.Type != tast.SymStruct {
			b.Errorf(expr.Sub.Pos, "symbol %s is a %s, not a struct",
				name, tast.SymStr(s.Type),
			)
			return nil
		}

		return s.Obj.(*objType).ref.TypeType()
	}

	b.Errorf(ast.ExprPos(expr), "expect a type")
	return nil
}

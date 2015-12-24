package g8

import (
	"e8vm.io/e8vm/g8/ast"
	"e8vm.io/e8vm/g8/types"
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

func buildPkgRef(b *builder, expr ast.Expr) *types.Pkg {
	switch expr := expr.(type) {
	case *ast.Operand:
		ret := buildOperand(b, expr)
		if ret == nil {
			return nil
		}
		if !ret.IsPkg() {
			b.Errorf(ast.ExprPos(expr), "expect a package, got %s", ret)
			return nil
		}
		return ret.Type().(*types.Pkg)
	}

	b.Errorf(ast.ExprPos(expr), "expect an imported package")
	return nil
}

func buildType(b *builder, expr ast.Expr) types.T {
	if expr == nil {
		panic("bug")
	}

	switch expr := expr.(type) {
	case *ast.Operand:
		ret := buildOperand(b, expr)
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
		pkg := buildPkgRef(b, expr.Expr)
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

		if s.Type != symStruct {
			b.Errorf(expr.Sub.Pos, "symbol %s is a %s, not a struct",
				name, symStr(s.Type),
			)
			return nil
		}

		return s.Obj.(*objType).ref.TypeType()
	}

	b.Errorf(ast.ExprPos(expr), "expect a type")
	return nil
}

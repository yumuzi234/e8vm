package sempass

import (
	"shanhu.io/smlvm/lexing"
	"shanhu.io/smlvm/pl/ast"
	"shanhu.io/smlvm/pl/tast"
	"shanhu.io/smlvm/pl/types"
	"shanhu.io/smlvm/syms"
)

const thisName = "<this>"

func buildFuncType(
	b *builder, recv *types.Pointer, f *ast.FuncSig,
) *types.Func {
	// the arguments
	args := buildParaList(b, f.Args)
	if args == nil {
		return nil
	}

	// the return values
	var rets []*types.Arg
	if f.RetType == nil {
		rets = buildParaList(b, f.Rets)
	} else {
		retType := buildType(b, f.RetType)
		if retType == nil {
			return nil
		}
		rets = []*types.Arg{{T: retType}}
	}

	if recv != nil {
		r := &types.Arg{Name: thisName, T: recv}
		return types.NewFunc(r, args, rets)
	}
	return types.NewFunc(nil, args, rets)
}

func buildArrayType(b *builder, expr *ast.ArrayTypeExpr) types.T {
	t := buildType(b, expr.Type)
	if t == nil {
		return nil
	}

	if expr.Len == nil {
		// slice
		return &types.Slice{T: t}
	}

	// array
	n := b.buildConstExpr(expr.Len)
	if n == nil {
		return nil
	}

	ntype := n.R().T
	ct, ok := ntype.(*types.Const)
	var v int64
	if !ok {
		// might be true, false, or other builtin consts
		b.CodeErrorf(ast.ExprPos(expr), "pl.nonConstArrayIndex",
			"array index is not a constant")
		return nil
	}

	if num, ok := types.NumConst(ntype); ok {
		v = num
	} else {
		v = ct.Value.(int64)
	}
	if v < 0 {
		b.CodeErrorf(ast.ExprPos(expr), "pl.negArrayIndex",
			"array index is negative: %d", v)
		return nil
	} else if !types.InRange(v, types.Int) {
		b.CodeErrorf(ast.ExprPos(expr), "pl.arrayIndexOutofRange",
			"index out of range of int32")
		return nil
	}
	return &types.Array{T: t, N: int32(v)}
}

func buildPkgRef(b *builder, ident *lexing.Token) *types.Pkg {
	s := b.scope.Query(ident.Lit)
	if s == nil {
		b.CodeErrorf(ident.Pos, "pl.undefinedIdent",
			"undefined identifier %s", ident.Lit)
		return nil
	}

	b.refSym(s, ident.Pos)
	if s.Type != tast.SymImport {
		b.Errorf(ident.Pos, "%s is not an imported package", ident.Lit)
		return nil
	}

	return s.ObjType.(*types.Pkg)
}

func buildType(b *builder, expr ast.Expr) types.T {
	if expr == nil {
		panic("bug")
	}
	hold := b.lhsSwap(false)
	defer b.lhsRestore(hold)

	switch expr := expr.(type) {
	case *ast.Operand:
		ret := buildOperand(b, expr)
		if ret == nil {
			return nil
		}
		ref := ret.R()
		t, ok := ref.T.(*types.Type)
		if !ok {
			b.CodeErrorf(ast.ExprPos(expr), "pl.expectType",
				"expect a type, got %s", ref.T)
			return nil
		}
		return t.T
	case *ast.StarExpr:
		t := buildType(b, expr.Expr)
		if t == nil {
			return nil
		}
		return &types.Pointer{T: t}
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
		if !syms.IsPublic(name) && s.Pkg() != b.path {
			b.Errorf(expr.Sub.Pos, "symbol %s is not public", name)
			return nil
		}

		if s.Type != tast.SymStruct {
			b.Errorf(expr.Sub.Pos, "symbol %s is a %s, not a struct",
				name, tast.SymStr(s.Type),
			)
			return nil
		}

		return s.ObjType.(*types.Type).T
	}

	b.CodeErrorf(ast.ExprPos(expr), "pl.expectType", "expect a type")
	return nil
}

package sempass

import (
	"e8vm.io/e8vm/g8/ast"
	"e8vm.io/e8vm/g8/tast"
	"e8vm.io/e8vm/g8/types"
	"e8vm.io/e8vm/lex8"
)

func buildIndexExpr(b *builder, expr *ast.IndexExpr) tast.Expr {
	array := b.BuildExpr(expr.Array)
	if array == nil {
		return nil
	}

	ref := array.R()
	if !ref.IsSingle() {
		b.Errorf(expr.Lbrack.Pos, "index on %s", ref)
		return nil
	}

	if expr.Colon != nil {
		return buildSlicing(b, expr, array)
	}

	return buildArrayGet(b, expr, array)
}

func elementType(t types.T) types.T {
	switch t := t.(type) {
	case *types.Array:
		return t.T
	case *types.Slice:
		return t.T
	}
	return nil
}

func checkArrayIndex(b *builder, index tast.Expr, pos *lex8.Pos) tast.Expr {
	t := index.R().T
	if v, ok := types.NumConst(t); ok {
		if v < 0 {
			b.Errorf(pos, "array index is negative: %d", v)
			return nil
		}
		return constCastInt(b, pos, v, index)
	}
	if !types.IsInteger(t) {
		b.Errorf(pos, "index must be an integer")
		return nil
	}
	return index
}

func buildArrayIndex(b *builder, expr ast.Expr, pos *lex8.Pos) tast.Expr {
	ret := b.BuildExpr(expr)
	if ret == nil {
		return nil
	}

	ref := ret.R()
	if !ref.IsSingle() {
		b.Errorf(pos, "index with %s", ref)
		return nil
	}
	return checkArrayIndex(b, ret, pos)
}

func buildSlicing(
	b *builder, expr *ast.IndexExpr, array tast.Expr,
) tast.Expr {
	t := array.R().T
	et := elementType(t)
	if et == nil {
		b.Errorf(expr.Lbrack.Pos, "slicing on neither array nor slice")
		return nil
	}

	var indexStart, indexEnd tast.Expr
	if expr.Index != nil {
		indexStart = buildArrayIndex(b, expr.Index, expr.Lbrack.Pos)
		if indexStart == nil {
			return nil
		}
	}

	if expr.IndexEnd != nil {
		indexEnd = buildArrayIndex(b, expr.IndexEnd, expr.Colon.Pos)
		if indexEnd == nil {
			return nil
		}
	}

	ref := tast.NewRef(&types.Slice{et})
	return &tast.IndexExpr{
		Array:    array,
		Index:    indexStart,
		IndexEnd: indexEnd,
		HasColon: true,
		Ref:      ref,
	}
}

func buildArrayGet(
	b *builder, expr *ast.IndexExpr, array tast.Expr,
) tast.Expr {
	t := array.R().T
	et := elementType(t)
	if et == nil {
		b.Errorf(expr.Lbrack.Pos, "index on neither array nor slice")
		return nil
	}

	index := buildArrayIndex(b, expr.Index, expr.Lbrack.Pos)
	if index == nil {
		return nil
	}

	ref := tast.NewAddressableRef(et)
	return &tast.IndexExpr{
		Array: array,
		Index: index,
		Ref:   ref,
	}
}

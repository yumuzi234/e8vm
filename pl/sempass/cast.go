package sempass

import (
	"shanhu.io/smlvm/pl/ast"
	"shanhu.io/smlvm/pl/tast"
	"shanhu.io/smlvm/pl/types"
)

func isPointerType(t types.T) bool {
	return types.IsPointer(t) || types.IsFuncPointer(t)
}

func regSizeCastable(to, from types.T) bool {
	if types.IsPointer(to) && types.IsPointer(from) {
		return true
	}
	if isPointerType(to) && types.IsBasic(from, types.Uint) {
		return true
	}
	if types.IsBasic(to, types.Uint) && isPointerType(from) {
		return true
	}
	return false
}

func buildCast(b *builder, expr *ast.CallExpr, t types.T) tast.Expr {
	pos := expr.Lparen.Pos

	args := buildExprList(b, expr.Args)
	if args == nil {
		return nil
	}

	ref := args.R()
	if !ref.IsSingle() {
		b.CodeErrorf(pos, "pl.cannotCast", "cannot convert %s to %s", ref, t)
		return nil
	}

	srcType := ref.T
	if c, ok := srcType.(*types.Const); ok {
		if v, ok := types.NumConst(srcType); ok && types.IsInteger(t) {
			return numCast(b, pos, v, args, t)
		}
		srcType = c.Type // using the underlying type
	}

	if types.IsInteger(t) && types.IsInteger(srcType) {
		return tast.NewCast(args, t)
	}
	if regSizeCastable(t, srcType) {
		return tast.NewCast(args, t)
	}
	if _, ok := t.(*types.Interface); ok {
		return tast.NewCast(args, t)
	}

	b.Errorf(pos, "cannot convert from %s to %s", srcType, t)
	return nil
}

func buildConstCast(b *builder, expr *ast.CallExpr, t types.T) tast.Expr {
	pos := expr.Lparen.Pos

	args := buildConstExprList(b, expr.Args)
	if args == nil {
		return nil
	}

	ref := args.R()
	if !ref.IsSingle() {
		b.CodeErrorf(pos, "pl.cannotCast", "cannot convert %s to %s", ref, t)
		return nil
	}

	srcType := args.Type().(*types.Const)
	ret := types.CastConst(srcType, t)
	if ret == nil {
		b.CodeErrorf(pos, "pl.cannotCast", "cannot convert %s to %s", ref, t)
		return nil
	}
	return tast.NewTypeConst(ret)
}

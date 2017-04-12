package sempass

import (
	"shanhu.io/smlvm/lexing"
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
		b.Errorf(pos, "cannot convert %s to %s", ref, t)
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

func implicitTypeCast(
	b *builder, pos *lexing.Pos, e tast.Expr, t types.T) tast.Expr {
	etype := e.Type()
	if types.IsNil(etype) {
		e = tast.NewCast(e, t)
	} else if v, ok := types.NumConst(etype); ok {
		e = numCast(b, pos, v, e, t)
		if e == nil {
			return nil
		}
	} else if _, ok := t.(*types.Interface); ok {
		e = tast.NewCast(e, t)
		if e == nil {
			panic("cannot cast interface")
		}
	}
	return e
}

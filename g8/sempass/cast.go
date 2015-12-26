package sempass

import (
	"e8vm.io/e8vm/g8/ast"
	"e8vm.io/e8vm/g8/tast"
	"e8vm.io/e8vm/g8/types"
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

func buildCast(b *Builder, expr *ast.CallExpr, t types.T) tast.Expr {
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
			return constCast(b, pos, v, ref, t)
		}
		srcType = c.Type // using the underlying type
	}

	if types.IsInteger(t) && types.IsInteger(srcType) {
		return &tast.Cast{args, tast.NewRef(t)}
	}
	if regSizeCastable(t, srcType) {
		return &tast.Cast{args, tast.NewRef(t)}
	}

	b.Errorf(pos, "cannot convert from %s to %s", srcType, t)
	return nil
}

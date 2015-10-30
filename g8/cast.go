package g8

import (
	"e8vm.io/e8vm/g8/ast"
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

func buildCast(b *builder, expr *ast.CallExpr, t types.T) *ref {
	pos := expr.Lparen.Pos

	args := buildExprList(b, expr.Args)
	if args == nil {
		return nil
	}

	if !args.IsSingle() {
		b.Errorf(pos, "cannot convert %s to %s", args, t)
		return nil
	}

	srcType := args.Type()
	ret := b.newTemp(t)
	if c, ok := srcType.(*types.Const); ok {
		if v, ok := types.NumConst(srcType); ok && types.IsInteger(t) {
			return constCast(b, pos, v, t)
		}
		srcType = c.Type // using the underlying type
	}

	if types.IsInteger(t) && types.IsInteger(srcType) {
		b.b.Arith(ret.IR(), nil, "cast", args.IR())
		return ret
	}
	if regSizeCastable(t, srcType) {
		b.b.Arith(ret.IR(), nil, "", args.IR())
		return ret
	}

	b.Errorf(pos, "cannot convert fromt %s to %s", srcType, t)
	return nil
}

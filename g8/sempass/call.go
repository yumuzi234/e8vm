package sempass

import (
	"e8vm.io/e8vm/fmt8"
	"e8vm.io/e8vm/g8/ast"
	"e8vm.io/e8vm/g8/tast"
	"e8vm.io/e8vm/g8/types"
)

func buildCallLen(b *Builder, expr *ast.CallExpr, f tast.Expr) tast.Expr {
	args := buildExprList(b, expr.Args)
	if args == nil {
		return nil
	}

	if args.Len() != 1 {
		b.Errorf(expr.Lparen.Pos, "len() takes one argument")
		return nil
	}

	ref := args.R()
	if !ref.IsSingle() {
		b.Errorf(expr.Lparen.Pos, "len() takes one argument")
		return nil
	}

	t := args.T
	switch t.(type) {
	case *types.Slice:
		return &tast.CallExpr{f, args, tast.NewRef(types.Int)}
	case *types.Array:
		return &tast.CallExpr{f, args, tast.NewRef(types.Int)}
	}

	b.Errorf(expr.Lparen.Pos, "len() does not take %s", t)
	return nil
}

func buildCallMake(b *Builder, expr *ast.CallExpr, f tast.Expr) tast.Expr {
	args := buildExprList(b, expr.Args)
	if args == nil {
		return nil
	}

	n := args.Len()
	if n == 0 {
		b.Errorf(expr.Lparen.Pos, "make() takes at least one argument")
		return nil
	}

	arg0 := args.Exprs[0]
	t, ok := arg0.R().T.(*types.Type)
	if !ok {
		b.Errorf(expr.Lparen.Pos, "make() takes a type as the 1st argument")
		return nil
	}
	switch t := t.T.(type) {
	case *types.Slice:
		if n != 3 {
			b.Errorf(expr.Lparen.Pos, "make() slice takes 3 arguments")
			return nil
		}

		size := args.Exprs[1]
		pos := ast.ExprPos(expr.Args.Exprs[1])
		size = checkArrayIndex(b, size, pos)
		if size == nil {
			return nil
		}

		start := args.Exprs[2]
		startType := start.R().T
		startPos := ast.ExprPos(expr.Args.Exprs[2])
		if v, ok := types.NumConst(startType); ok {
			start = constCastUint(b, startPos, v, start)
			if start == nil {
				return nil
			}
		} else if !types.IsBasic(startType, types.Uint) {
			pt := types.PointerOf(startType)
			if pt == nil || !types.SameType(pt, t.T) {
				b.Errorf(startPos,
					"make() takes an uint or a typed pointer as the 3rd arg",
				)
				return nil
			}
		}

		callArgs := tast.NewExprList()
		callArgs.Append(arg0)
		callArgs.Append(size)
		callArgs.Append(start)

		r := tast.NewAddressableRef(t.T)
		return &tast.CallExpr{f, callArgs, r}
	}

	b.Errorf(expr.Lparen.Pos, "cannot make() type %s", t.T)
	return nil
}

func buildCallExpr(b *Builder, expr *ast.CallExpr) tast.Expr {
	f := b.BuildExpr(expr.Func)
	if f == nil {
		return nil
	}

	pos := ast.ExprPos(expr.Func)
	fref := f.R()
	if !fref.IsSingle() {
		b.Errorf(pos, "%s is not callable", fref)
		return nil
	}
	if t, ok := fref.T.(*types.Type); ok {
		return buildCast(b, expr, t.T)
	}

	builtin, ok := fref.T.(*types.BuiltInFunc)
	if ok {
		switch builtin.Name {
		case "len":
			return buildCallLen(b, expr, f)
		case "make":
			return buildCallMake(b, expr, f)
		}
		b.Errorf(pos, "builtin %s() not implemented", builtin.Name)
		return nil
	}

	funcType, ok := fref.T.(*types.Func)
	if !ok {
		b.Errorf(pos, "function call on non-callable: %s", fref)
		return nil
	}

	args := buildExprList(b, expr.Args)
	if args == nil {
		return nil
	}

	nargs := args.Len()
	if nargs != len(funcType.Args) {
		b.Errorf(ast.ExprPos(expr), "argument expects (%s), got (%s)",
			fmt8.Join(funcType.Args, ","), args,
		)
		return nil
	}

	callArgs := tast.NewExprList()

	// type check and perform casting on each argument
	for i := 0; i < nargs; i++ {
		argExpr := args.Exprs[i]
		argType := argExpr.Type()
		expect := funcType.Args[i].T
		if !types.CanAssign(expect, argType) {
			pos := ast.ExprPos(expr.Args.Exprs[i])
			b.Errorf(pos, "argument %d expects %s, got %s",
				i+1, expect, argType,
			)
			return nil
		}

		// insert auto casts
		if types.IsNil(argType) {
			callArgs.Append(tast.NewCast(argExpr, expect))
		} else if _, ok := types.NumConst(argType); ok {
			callArgs.Append(tast.NewCast(argExpr, expect))
		} else {
			callArgs.Append(argExpr)
		}
	}

	retRef := tast.Void
	for _, t := range funcType.RetTypes {
		retRef = tast.AppendRef(retRef, tast.NewRef(t))
	}

	return &tast.CallExpr{f, callArgs, retRef}
}

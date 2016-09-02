package sempass

import (
	"e8vm.io/e8vm/fmtutil"
	"e8vm.io/e8vm/pl/ast"
	"e8vm.io/e8vm/pl/tast"
	"e8vm.io/e8vm/pl/types"
)

func buildCallLen(b *builder, expr *ast.CallExpr, f tast.Expr) tast.Expr {
	args := buildExprList(b, expr.Args)
	if args == nil {
		return nil
	}

	ref := args.R()
	if !ref.IsSingle() {
		b.Errorf(expr.Lparen.Pos, "len() takes one argument")
		return nil
	}

	t := ref.T
	switch t.(type) {
	case *types.Slice:
		return &tast.CallExpr{f, args, tast.NewRef(types.Int)}
	case *types.Array:
		return &tast.CallExpr{f, args, tast.NewRef(types.Int)}
	}

	b.Errorf(expr.Lparen.Pos, "len() does not take %s", t)
	return nil
}

func buildCallMake(b *builder, expr *ast.CallExpr, f tast.Expr) tast.Expr {
	args := buildExprList(b, expr.Args)
	if args == nil {
		return nil
	}

	n := args.R().Len()
	if n == 0 {
		b.Errorf(expr.Lparen.Pos, "make() takes at least one argument")
		return nil
	}

	argsList, ok := tast.MakeExprList(args)
	if !ok {
		b.Errorf(expr.Lparen.Pos, "make() only takes a literal list")
		return nil
	}

	arg0 := argsList.Exprs[0]
	t, ok := arg0.R().T.(*types.Type)
	if !ok {
		b.Errorf(expr.Lparen.Pos, "make() takes a type as the 1st argument")
		return nil
	}
	switch st := t.T.(type) {
	case *types.Slice:
		if n != 3 {
			b.Errorf(expr.Lparen.Pos, "make() slice takes 3 arguments")
			return nil
		}

		size := argsList.Exprs[1]
		pos := ast.ExprPos(expr.Args.Exprs[1])
		size = checkArrayIndex(b, size, pos)
		if size == nil {
			return nil
		}

		start := argsList.Exprs[2]
		startType := start.R().T
		startPos := ast.ExprPos(expr.Args.Exprs[2])
		if v, ok := types.NumConst(startType); ok {
			start = constCastUint(b, startPos, v, start)
			if start == nil {
				return nil
			}
		} else if !types.IsBasic(startType, types.Uint) {
			pt := types.PointerOf(startType)
			if pt == nil || !types.SameType(pt, st.T) {
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

		r := tast.NewRef(st)
		return &tast.CallExpr{f, callArgs, r}
	}

	b.Errorf(expr.Lparen.Pos, "cannot make() type %s", t.T)
	return nil
}

func buildCallExpr(b *builder, expr *ast.CallExpr) tast.Expr {
	hold := b.lhsSwap(false)
	defer b.lhsRestore(hold)

	f := b.buildExpr(expr.Func)
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

	argsRef := args.R()
	nargs := argsRef.Len()
	if nargs != len(funcType.Args) {
		b.Errorf(ast.ExprPos(expr), "argument expects (%s), got (%s)",
			fmtutil.Join(funcType.Args, ","), args,
		)
		return nil
	}

	// type check on each argument
	for i := 0; i < nargs; i++ {
		argType := argsRef.At(i).Type()
		expect := funcType.Args[i].T
		if !types.CanAssign(expect, argType) {
			pos := ast.ExprPos(expr.Args.Exprs[i])
			b.Errorf(pos, "argument %d expects %s, got %s",
				i+1, expect, argType,
			)
			return nil
		}
	}

	// insert casting when it is a literal expression list.
	callArgs, ok := tast.MakeExprList(args)
	if ok {
		castedArgs := tast.NewExprList()
		for i := 0; i < nargs; i++ {
			argExpr := callArgs.Exprs[i]
			argType := argExpr.Type()
			expect := funcType.Args[i].T

			// insert auto casts for consts
			if types.IsNil(argType) {
				castedArgs.Append(tast.NewCast(argExpr, expect))
			} else if _, ok := types.NumConst(argType); ok {
				castedArgs.Append(tast.NewCast(argExpr, expect))
			} else {
				castedArgs.Append(argExpr)
			}
		}
		args = castedArgs
	}

	retRef := tast.Void
	for _, t := range funcType.RetTypes {
		retRef = tast.AppendRef(retRef, tast.NewRef(t))
	}

	return &tast.CallExpr{f, args, retRef}
}

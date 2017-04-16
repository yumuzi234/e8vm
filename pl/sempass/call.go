package sempass

import (
	"shanhu.io/smlvm/fmtutil"
	"shanhu.io/smlvm/pl/ast"
	"shanhu.io/smlvm/pl/tast"
	"shanhu.io/smlvm/pl/types"
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
		return &tast.CallExpr{Func: f, Args: args, Ref: tast.NewRef(types.Int)}
	case *types.Array:
		return &tast.CallExpr{Func: f, Args: args, Ref: tast.NewRef(types.Int)}
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
			start = numCastUint(b, startPos, v, start)
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
		return &tast.CallExpr{Func: f, Args: callArgs, Ref: r}
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
		b.CodeErrorf(pos, "pl.call.notSingle", "%s is not callable", fref)
		return nil
	}

	// expr.Func is a Type

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
	pos = ast.ExprPos(expr)
	if nargs != len(funcType.Args) {
		b.CodeErrorf(pos, "pl.argsMismatch.count",
			"argument count mismatch, expects (%s), got (%s)",
			fmtutil.Join(funcType.Args, ","), args,
		)
		return nil
	}

	srcTypes := argsRef.TypeList()
	destTypes := funcType.ArgTypes
	ok, needCast, castMask := canAssigns(b, pos, destTypes, srcTypes)
	if !ok {
		return nil
	}
	if needCast {
		args = tast.NewMultiCastTypes(args, destTypes, castMask)
	}

	retRef := tast.NewListRef(funcType.RetTypes)
	return &tast.CallExpr{Func: f, Args: args, Ref: retRef}
}

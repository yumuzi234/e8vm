package g8

import (
	"e8vm.io/e8vm/fmt8"
	"e8vm.io/e8vm/g8/ast"
	"e8vm.io/e8vm/g8/ir"
	"e8vm.io/e8vm/g8/types"
)

func buildCallMake(b *builder, expr *ast.CallExpr) *ref {
	args := buildExprList(b, expr.Args)
	if args == nil {
		return nil
	}

	n := args.Len()
	if n == 0 {
		b.Errorf(expr.Lparen.Pos, "make() takes at least one argument")
		return nil
	}

	arg0 := args.At(0)
	if !arg0.IsType() {
		b.Errorf(expr.Lparen.Pos, "make() takes a type as the 1st argument")
		return nil
	}

	t := arg0.Type().(*types.Type).T
	switch t := t.(type) {
	case *types.Slice:
		if n != 3 {
			b.Errorf(expr.Lparen.Pos, "make() slice takes 3 arguments")
			return nil
		}

		size := args.At(1)
		pos := ast.ExprPos(expr.Args.Exprs[1])
		sizeIr := checkArrayIndex(b, size, pos)
		if sizeIr == nil {
			return nil
		}

		start := args.At(2)
		startType := start.Type()
		startPos := ast.ExprPos(expr.Args.Exprs[2])
		if v, ok := types.NumConst(startType); ok {
			start = constCastUint(b, startPos, v)
			if start == nil {
				return nil
			}
		} else if !types.IsBasic(startType, types.Uint) {
			pt := types.PointerOf(startType)
			if pt == nil || !types.SameType(pt, t.T) {
				b.Errorf(startPos,
					"make() takes an uint or a typed pointer as 3rd argument",
				)
				return nil
			}
		}

		return newSlice(b, t.T, start.IR(), sizeIr)
	}

	b.Errorf(expr.Lparen.Pos, "cannot make() type %s", t)
	return nil
}

func buildCallLen(b *builder, expr *ast.CallExpr) *ref {
	args := buildExprList(b, expr.Args)
	if args == nil {
		return nil
	}

	if !args.IsSingle() {
		b.Errorf(expr.Lparen.Pos, "len() takes one argument")
		return nil
	}

	t := args.Type()
	ret := b.newTemp(types.Int)

	switch t := t.(type) {
	case *types.Slice:
		addr := b.newPtr()
		b.b.Arith(addr, nil, "&", args.IR())
		b.b.Assign(ret.IR(), ir.NewAddrRef(addr, 4, 4, false, true))
	case *types.Array:
		b.b.Assign(ret.IR(), ir.Num(uint32(t.N)))
	default:
		b.Errorf(expr.Lparen.Pos, "len() does not take %s", t)
		return nil
	}

	return ret
}

func buildCallExpr(b *builder, expr *ast.CallExpr) *ref {
	f := b.buildExpr(expr.Func)
	if f == nil {
		return nil
	}

	pos := ast.ExprPos(expr.Func)

	if !f.IsSingle() {
		b.Errorf(pos, "expression list is not callable")
		return nil
	} else if f.IsType() {
		return buildCast(b, expr, f.TypeType())
	}

	builtin, ok := f.Type().(*types.BuiltInFunc)
	if ok {
		switch builtin.Name {
		case "len":
			return buildCallLen(b, expr)
		case "make":
			return buildCallMake(b, expr)
		}
		b.Errorf(pos, "builtin %s() not implemented", builtin.Name)
		return nil
	}

	nilFuncPointerPanic(b, f.IR())

	funcType, ok := f.Type().(*types.Func)
	if !ok {
		// not a function
		b.Errorf(pos, "function call on non-callable: %s", f)
		return nil
	}

	args := buildExprList(b, expr.Args)
	if args == nil {
		return nil
	}

	nargs := args.Len()
	if nargs != len(funcType.Args) {
		b.Errorf(ast.ExprPos(expr), "argument expects (%s), got (%s)",
			fmt8.Join(funcType.Args, ","), fmt8.Join(args.TypeList(), ","),
		)
		return nil
	}

	// type check on parameters
	for i := 0; i < nargs; i++ {
		argType := args.At(i).Type()
		expect := funcType.Args[i].T
		if !types.CanAssign(expect, argType) {
			pos := ast.ExprPos(expr.Args.Exprs[i])
			b.Errorf(pos, "argument %d expects %s, got %s",
				i+1, expect, argType,
			)
			return nil
		}
	}

	argsCasted := new(ref)

	// auto type casts for nil and consts.
	for i := 0; i < nargs; i++ {
		argRef := args.At(i)
		argType := argRef.Type()
		expect := funcType.Args[i].T

		if types.IsNil(argType) {
			tmp := b.newTemp(expect).IR()
			b.b.Zero(tmp)
			argsCasted = appendRef(argsCasted, newRef(expect, tmp))
		} else if v, ok := types.NumConst(argType); ok {
			tmp := b.newTemp(expect).IR()
			b.b.Assign(tmp, constNumIr(v, expect))
			argsCasted = appendRef(argsCasted, newRef(expect, tmp))
		} else {
			argsCasted = appendRef(argsCasted, argRef)
		}
	}

	ret := new(ref)
	for _, t := range funcType.RetTypes {
		ret = appendRef(ret, newRef(t, b.newTempIR(t)))
	}

	// call the func in IR
	if f.recv == nil {
		irs := argsCasted.IRList()
		fref := wrapFuncPtr(f.IR(), funcType)
		b.b.Call(ret.IRList(), fref, irs...)
	} else {
		var irs []ir.Ref
		irs = append(irs, f.recv.IR())
		irs = append(irs, argsCasted.IRList()...)
		fref := wrapFuncPtr(f.IR(), f.recvFunc)
		b.b.Call(ret.IRList(), fref, irs...)
	}

	return ret
}

func wrapFuncPtr(f ir.Ref, t *types.Func) ir.Ref {
	switch f := f.(type) {
	case *ir.FuncSym:
		return f
	case *ir.Func:
		return f
	}
	return ir.NewFuncPtr(makeFuncSig(t), f)
}

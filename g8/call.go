package g8

import (
	"e8vm.io/e8vm/fmt8"
	"e8vm.io/e8vm/g8/ast"
	"e8vm.io/e8vm/g8/ir"
	"e8vm.io/e8vm/g8/types"
)

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
		}
		b.Errorf(pos, "builtin %s() not implemented", builtin.Name)
		return nil
	}

	nilFuncPointerPanic(b, f.IR())

	funcType, ok := f.Type().(*types.Func) // the func sig in the builder
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
			argsCasted = appendRef(argsCasted, newRef(argType, tmp))
		} else if v, ok := types.NumConst(argType); ok {
			tmp := b.newTemp(expect).IR()
			b.b.Assign(tmp, constNumIr(v, expect))
			argsCasted = appendRef(argsCasted, newRef(argType, tmp))
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
		b.b.Call(ret.IRList(), f.IR(), funcType.Sig, irs...)
	} else {
		var irs []ir.Ref
		irs = append(irs, f.recv.IR())
		irs = append(irs, argsCasted.IRList()...)

		b.b.Call(ret.IRList(), f.IR(), f.recvFunc.Sig, irs...)
	}

	return ret
}

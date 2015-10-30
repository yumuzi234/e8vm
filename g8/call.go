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

	if args.Len() != len(funcType.Args) {
		b.Errorf(ast.ExprPos(expr), "argument expects (%s), got (%s)",
			fmt8.Join(funcType.Args, ","), fmt8.Join(args.typ, ","),
		)
		return nil
	}

	// type check on parameters
	for i, argType := range args.typ {
		expect := funcType.Args[i].T
		if !types.CanAssign(expect, argType) {
			pos := ast.ExprPos(expr.Args.Exprs[i])
			b.Errorf(pos, "argument %d expects %s, got %s",
				i+1, expect, argType,
			)
			return nil
		}
	}

	for i, argType := range args.typ {
		expect := funcType.Args[i].T
		if types.IsNil(argType) {
			tmp := b.newTemp(expect).IR()
			b.b.Zero(tmp)
			args.ir[i] = tmp
		} else if v, ok := types.NumConst(argType); ok {
			tmp := b.newTemp(expect).IR()
			b.b.Assign(tmp, constNumIr(v, expect))
			args.ir[i] = tmp
		}
	}

	ret := new(ref)
	ret.typ = funcType.RetTypes
	for _, t := range funcType.RetTypes {
		ret.ir = append(ret.ir, b.newTempIR(t))
		ret.addressable = append(ret.addressable, false)
	}

	// call the func in IR
	if f.recv == nil {
		b.b.Call(ret.ir, f.IR(), funcType.Sig, args.ir...)
	} else {
		var irs []ir.Ref
		irs = append(irs, f.recv.IR())
		irs = append(irs, args.ir...)

		b.b.Call(ret.ir, f.IR(), f.recvFunc.Sig, irs...)
	}

	return ret
}

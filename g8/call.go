package g8

import (
	"e8vm.io/e8vm/g8/ir"
	"e8vm.io/e8vm/g8/tast"
	"e8vm.io/e8vm/g8/types"
)

func genCallLen(b *builder, expr *tast.CallExpr) *ref {
	args := genExprList(b, expr.Args)
	t := args.Type()
	ret := b.newTemp(types.Int)
	switch t := t.(type) {
	case *types.Slice:
		addr := b.newPtr()
		b.b.Arith(addr, nil, "&", args.IR())
		b.b.Assign(ret.IR(), ir.NewAddrRef(addr, 4, 4, false, true))
		return ret
	case *types.Array:
		b.b.Assign(ret.IR(), ir.Num(uint32(t.N)))
		return ret
	}
	panic("bug")
}

func genCallMake(b *builder, expr *tast.CallExpr) *ref {
	args := genExprList(b, expr.Args)
	arg0 := args.At(0)
	t := arg0.Type().(*types.Type).T.(*types.Slice)
	size := checkArrayIndex2(b, args.At(1))
	start := args.At(2).IR()
	return newSlice(b, t.T, start, size)
}

func genCallExpr(b *builder, expr *tast.CallExpr) *ref {
	f := b.buildExpr2(expr.Func)
	builtin, ok := f.Type().(*types.BuiltInFunc)
	if ok {
		switch builtin.Name {
		case "len":
			return genCallLen(b, expr)
		case "make":
			return genCallMake(b, expr)
		}
		panic("bug")
	}

	nilFuncPointerPanic(b, f.IR())
	funcType := f.Type().(*types.Func)
	args := genExprList(b, expr.Args)

	ret := new(ref)
	for _, t := range funcType.RetTypes {
		ret = appendRef(ret, newRef(t, b.newTempIR(t)))
	}

	if f.recv == nil {
		irs := args.IRList()
		fref := wrapFuncPtr(f.IR(), funcType)
		b.b.Call(ret.IRList(), fref, irs...)
	} else {
		var irs []ir.Ref
		irs = append(irs, f.recv.IR())
		irs = append(irs, args.IRList()...)
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

package pl

import (
	"shanhu.io/smlvm/pl/codegen"
	"shanhu.io/smlvm/pl/tast"
	"shanhu.io/smlvm/pl/types"
)

func buildCallLen(b *builder, expr *tast.CallExpr) *ref {
	args := b.buildExpr(expr.Args)
	t := args.Type()
	ret := b.newTemp(types.Int)
	switch t := t.(type) {
	case *types.Slice:
		addr := b.newPtr()
		b.b.Arith(addr, nil, "&", args.IR())
		b.b.Assign(ret.IR(), codegen.NewAddrRef(addr, 4, 4, false, true))
		return ret
	case *types.Array:
		b.b.Assign(ret.IR(), codegen.Num(uint32(t.N)))
		return ret
	}
	panic("bug")
}

func buildCallMake(b *builder, expr *tast.CallExpr) *ref {
	args := b.buildExpr(expr.Args)
	arg0 := args.At(0)
	t := arg0.Type().(*types.Type).T.(*types.Slice)
	size := checkArrayIndex(b, args.At(1))
	start := args.At(2).IR()
	nilPointerPanic(b, start)
	return newSlice(b, t.T, start, size)
}

func buildCallExpr(b *builder, expr *tast.CallExpr) *ref {
	f := b.buildExpr(expr.Func)
	builtin, ok := f.Type().(*types.BuiltInFunc)
	if ok {
		switch builtin.Name {
		case "len":
			return buildCallLen(b, expr)
		case "make":
			return buildCallMake(b, expr)
		}
		panic("bug")
	}

	nilFuncPointerPanic(b, f.IR())
	funcType := f.Type().(*types.Func)
	args := b.buildExpr(expr.Args)

	ret := new(ref)
	for _, t := range funcType.RetTypes {
		ret = appendRef(ret, newRef(t, b.newTempIR(t)))
	}

	if f.recv == nil {
		irs := args.IRList()
		fref := wrapFuncPtr(f.IR(), funcType)
		b.b.Call(ret.IRList(), fref, irs...)
	} else {
		var irs []codegen.Ref
		irs = append(irs, f.recv.IR())
		irs = append(irs, args.IRList()...)
		fref := wrapFuncPtr(f.IR(), f.recvFunc)
		b.b.Call(ret.IRList(), fref, irs...)
	}

	return ret
}

func wrapFuncPtr(f codegen.Ref, t *types.Func) codegen.Ref {
	switch f := f.(type) {
	case *codegen.FuncSym:
		return f
	case *codegen.Func:
		return f
	}
	return codegen.NewFuncPtr(makeFuncSig(t), f)
}

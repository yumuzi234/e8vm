package g8

import (
	"e8vm.io/e8vm/arch8"
	"e8vm.io/e8vm/g8/ir"
	"e8vm.io/e8vm/g8/types"
)

const (
	startName     = ":start"
	testStartName = ":test"
)

func findFunc(b *builder, name string, t types.T) *objFunc {
	s := b.scope.Query(name)
	if s == nil {
		return nil
	}
	f, isFunc := s.Item.(*objFunc)
	if !isFunc {
		return nil
	}
	if !types.SameType(f.ref.Type(), t) {
		return nil
	}
	return f
}

func addStart(b *builder) {
	mainFunc := findFunc(b, "main", types.VoidFunc)
	if mainFunc == nil {
		return
	}

	b.f = b.p.NewFunc(startName, nil, ir.VoidFuncSig)
	b.f.SetAsMain()
	b.b = b.f.NewBlock(nil)
	b.b.Call(nil, mainFunc.IR(), ir.VoidFuncSig)
}

var testMainFuncType = types.NewVoidFunc(types.VoidFunc)
var testMainFuncSig = ir.NewFuncSig(
	[]*ir.FuncArg{{
		Name:         "f",
		Size:         arch8.RegSize,
		U8:           false,
		RegSizeAlign: true,
	}},
	nil,
)

func addTestStart(b *builder, testList ir.Ref, n int) {
	b.f = b.p.NewFunc(testStartName, nil, ir.VoidFuncSig)
	b.f.SetAsMain()
	b.b = b.f.NewBlock(nil)

	argAddr := ir.Num(arch8.AddrBootArg) // the arg
	index := b.newTempIR(types.Uint)     // to save the index
	b.b.Assign(index, ir.NewAddrRef(argAddr, arch8.RegSize, 0, false, true))

	size := ir.Num(uint32(n))
	checkInRange(b, index, size, "u<")

	base := b.newPtr()
	b.b.Arith(base, nil, "&", testList)
	addr := b.newPtr()
	b.b.Arith(addr, index, "*", ir.Num(arch8.RegSize))
	b.b.Arith(addr, base, "+", addr)

	f := ir.NewAddrRef(addr, arch8.RegSize, 0, false, true)

	testMain := findFunc(b, "testMain", testMainFuncType)
	if testMain == nil {
		b.b.Call(nil, f, ir.VoidFuncSig)
	} else {
		b.b.Call(nil, testMain.ref.IR(), testMainFuncSig, f)
	}
}

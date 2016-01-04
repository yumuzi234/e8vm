package g8

import (
	"e8vm.io/e8vm/g8/ir"
	"e8vm.io/e8vm/g8/types"
)

func makeArg(t *types.Arg) *ir.FuncArg {
	return &ir.FuncArg{
		Name:         t.Name,
		Size:         t.Size(),
		U8:           types.IsBasic(t.T, types.Uint8),
		RegSizeAlign: t.RegSizeAlign(),
	}
}

// converts a langauge function signature into a IR function signature
func makeFuncSig(f *types.Func) *ir.FuncSig {
	narg := len(f.Args)
	args := make([]*ir.FuncArg, 0, narg)

	for _, t := range f.Args {
		if t.T == nil {
			panic("type missing")
		}
		args = append(args, makeArg(t))
	}

	rets := make([]*ir.FuncArg, len(f.Rets))
	for i, t := range f.Rets {
		rets[i] = makeArg(t)
	}

	return ir.NewFuncSig(args, rets)
}

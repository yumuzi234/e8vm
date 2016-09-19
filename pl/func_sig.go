package pl

import (
	"shanhu.io/smlvm/pl/codegen"
	"shanhu.io/smlvm/pl/types"
)

func makeArg(t *types.Arg) *codegen.FuncArg {
	return &codegen.FuncArg{
		Name:         t.Name,
		Size:         t.Size(),
		U8:           types.IsBasic(t.T, types.Uint8),
		RegSizeAlign: t.RegSizeAlign(),
	}
}

// converts a langauge function signature into a IR function signature
func makeFuncSig(f *types.Func) *codegen.FuncSig {
	narg := len(f.Args)
	args := make([]*codegen.FuncArg, 0, narg)

	for _, t := range f.Args {
		if t.T == nil {
			panic("type missing")
		}
		args = append(args, makeArg(t))
	}

	rets := make([]*codegen.FuncArg, len(f.Rets))
	for i, t := range f.Rets {
		rets[i] = makeArg(t)
	}

	return codegen.NewFuncSig(args, rets)
}

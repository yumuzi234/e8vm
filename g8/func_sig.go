package g8

import (
	"e8vm.io/e8vm/g8/ast"
	"e8vm.io/e8vm/g8/types"
)

const thisName = "<this>"

func buildFuncType(b *builder, s *structInfo, f *ast.FuncSig) *types.Func {
	// the arguments
	args := buildParaList(b, f.Args)
	if args == nil {
		return nil
	}

	// the return values
	var rets []*types.Arg
	if f.RetType == nil {
		rets = buildParaList(b, f.Rets)
	} else {
		retType := b.buildType(f.RetType)
		if retType == nil {
			return nil
		}
		rets = []*types.Arg{{T: retType}}
	}

	if s != nil {
		recv := &types.Arg{Name: thisName, T: s.pt}
		return types.NewFunc(recv, args, rets)
	}

	return types.NewFunc(nil, args, rets)
}

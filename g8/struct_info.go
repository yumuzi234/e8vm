package g8

import (
	"e8vm.io/e8vm/g8/types"
	"e8vm.io/e8vm/lex8"
)

type structInfo struct {
	name *lex8.Token
	t    *types.Struct // the struct type
	pt   *types.Pointer
}

func (info *structInfo) Name() string {
	return info.name.Lit
}

package codegen

import (
	"e8vm.io/e8vm/arch"
)

const regSize = arch.RegSize

func alignUp(size, align int32) int32 {
	mod := size % align
	if mod == 0 {
		return size
	}
	return size + align - mod
}

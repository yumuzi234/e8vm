package test

import (
	"e8vm.io/e8vm/glang"
)

// Fuzz implements go-fuzz interface.
func Fuzz(data []byte) int {
	_, errs, _ := glang.CompileSingle("test.g", string(data), false)
	if len(errs) > 0 {
		return 0
	}
	return 1
}

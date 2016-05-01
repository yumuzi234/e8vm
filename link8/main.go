package link8

import (
	"e8vm.io/e8vm/arch8"
	"e8vm.io/e8vm/asm8/inst"
)

func wrapMain(funcs []*PkgSym) *Func {
	ret := NewFunc() // the main func

	// clear r0 for safety
	ret.AddInst(inst.Reg(arch8.XOR, 0, 0, 0, 0, 0))

	for _, f := range funcs {
		ret.AddInst(inst.Jmp(arch8.JAL, 0))
		ret.AddLink(FillLink, f)
	}

	ret.AddInst(inst.Sys(arch8.HALT, 0, 0))

	return ret
}

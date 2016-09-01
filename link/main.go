package link

import (
	"e8vm.io/e8vm/arch"
	"e8vm.io/e8vm/asm/inst"
)

func wrapMain(funcs []*PkgSym) *Func {
	ret := NewFunc() // the main func

	// clear r0 for safety
	ret.AddInst(inst.Reg(arch.XOR, 0, 0, 0, 0, 0))

	for _, f := range funcs {
		ret.AddInst(inst.Jmp(arch.JAL, 0))
		ret.AddLink(FillLink, f)
	}

	ret.AddInst(inst.Sys(arch.HALT, 0, 0))

	return ret
}

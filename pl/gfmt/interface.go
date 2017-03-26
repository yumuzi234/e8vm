package gfmt

import (
	"shanhu.io/smlvm/pl/ast"
)

func printInterface(f *formatter, d *ast.Interface) {
	f.printExprs(d.Kw, " ", d.Name, " ", d.Lbrace)
	if len(d.Funcs) == 0 {
		f.printToken(d.Rbrace)
		return
	}
	f.printEndl()
	f.Tab()
	for i, fun := range d.Funcs {
		if i != 0 {
			f.printGap()
		}
		f.printExprs(fun.Name, " ")
		f.printSpace()
		printFuncSig(f, fun.FuncSigs)
	}
	f.printEndl()
	f.ShiftTab()
	f.printToken(d.Rbrace)
}

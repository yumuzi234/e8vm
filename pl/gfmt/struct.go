package gfmt

import (
	"shanhu.io/smlvm/pl/ast"
)

func printStruct(f *formatter, d *ast.Struct) {
	f.printExprs(d.Kw, " ", d.Name, " ", d.Lbrace)
	if len(d.Fields) == 0 {
		f.printToken(d.Rbrace)
		return
	}
	f.printEndl()
	f.Tab()
	for i, field := range d.Fields {
		if i > 0 {
			f.printGap()
		}
		printIdents(f, field.Idents)
		f.printSpace()
		f.printExprs(field.Type)
	}
	f.printEndl()
	f.ShiftTab()
	f.printToken(d.Rbrace)
}

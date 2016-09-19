package gfmt

import (
	"shanhu.io/smlvm/pl/ast"
)

func printStruct(f *formatter, d *ast.Struct) {
	f.printExprs(d.Kw, " ", d.Name, " ", d.Lbrace)
	f.printEndl()
	f.Tab()

	for i, field := range d.Fields {
		printIdents(f, field.Idents)
		f.printSpace()
		f.printExprs(field.Type)
		f.printEndlPlus(i < len(d.Fields)-1, false)
	}

	if len(d.Methods) > 0 {
		f.printEndl()
	}
	for i, method := range d.Methods {
		printFunc(f, method)
		f.printEndlPlus(i < len(d.Methods)-1, false)
	}

	f.ShiftTab()
	f.printToken(d.Rbrace)
}

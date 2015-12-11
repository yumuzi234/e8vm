package gfmt

import (
	"e8vm.io/e8vm/g8/ast"
)

func printStruct(f *formatter, d *ast.Struct) {
	printExprs(f, d.Kw, " ", d.Name, " ", d.Lbrace)
	f.printEndl()
	f.Tab()

	for _, field := range d.Fields {
		printIdents(f, field.Idents)
		f.printSpace()
		printExpr(f, field.Type)
		f.printEndl()
	}

	f.printEndl()
	for _, method := range d.Methods {
		printFunc(f, method)
	}

	f.ShiftTab()
	f.printToken(d.Rbrace)
	f.printEndl()
}

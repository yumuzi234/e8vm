package gfmt

import (
	"fmt"

	"e8vm.io/e8vm/fmt8"
	"e8vm.io/e8vm/g8/ast"
)

func printStruct(p *fmt8.Printer, m *matcher, d *ast.Struct) {
	printExprs(p, m, d.Kw, " ", d.Name, " ", d.Lbrace)
	fmt.Fprintln(p)
	p.Tab()

	for _, field := range d.Fields {
		printIdents(p, m, field.Idents)
		fmt.Fprint(p, " ")
		printExpr(p, m, field.Type)
		fmt.Fprintln(p)
	}

	fmt.Fprintln(p)
	for _, method := range d.Methods {
		printFunc(p, m, method)
	}

	p.ShiftTab()
	printToken(p, m, d.Rbrace)
	fmt.Fprintln(p)
}

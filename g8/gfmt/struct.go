package gfmt

import (
	"fmt"

	"e8vm.io/e8vm/fmt8"
	"e8vm.io/e8vm/g8/ast"
)

func printStruct(p *fmt8.Printer, d *ast.Struct) {
	fmt.Fprintf(p, "struct %s {\n", d.Name.Lit)
	p.Tab()

	for _, field := range d.Fields {
		printIdents(p, field.Idents)
		fmt.Fprint(p, " ")
		printExpr(p, field.Type)
		fmt.Fprintln(p)
	}

	fmt.Fprintln(p)
	for _, method := range d.Methods {
		printFunc(p, method)
	}

	p.ShiftTab()
	fmt.Fprintln(p, "}")
}

package ast

import (
	"fmt"

	"e8vm.io/e8vm/fmt8"
)

func printIdents(p *fmt8.Printer, idents *IdentList) {
	ss := make([]string, len(idents.Idents))
	for i, id := range idents.Idents {
		ss[i] = id.Lit
	}
	fmt.Fprint(p, fmt8.Join(ss, ", "))
}

func printVarDecl(p *fmt8.Printer, d *VarDecl) {
	printIdents(p, d.Idents)
	if d.Type != nil {
		fmt.Fprint(p, " ")
		printExpr(p, d.Type)
	}

	if d.Eq != nil {
		fmt.Fprint(p, " = ")
		printExpr(p, d.Exprs)
	}
}

func printVarDecls(p *fmt8.Printer, d *VarDecls) {
	if d.Lparen == nil {
		// single declare
		fmt.Fprintf(p, "var ")
		for _, decl := range d.Decls {
			printVarDecl(p, decl)
		}
	} else {
		fmt.Fprintf(p, "var (\n")
		p.Tab()
		for _, decl := range d.Decls {
			printVarDecl(p, decl)
			fmt.Println(p)
		}
		p.ShiftTab()
		fmt.Fprintf(p, ")")
	}
}

func printConstDecls(p *fmt8.Printer, d *ConstDecls) {
	fmt.Fprintf(p, "<todo: const decls>")
}

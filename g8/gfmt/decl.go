package gfmt

import (
	"fmt"

	"e8vm.io/e8vm/fmt8"
	"e8vm.io/e8vm/g8/ast"
)

func printIdents(p *fmt8.Printer, idents *ast.IdentList) {
	ss := make([]string, len(idents.Idents))
	for i, id := range idents.Idents {
		ss[i] = id.Lit
	}
	fmt.Fprint(p, fmt8.Join(ss, ", "))
}

func printVarDecl(p *fmt8.Printer, d *ast.VarDecl) {
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

func printVarDecls(p *fmt8.Printer, d *ast.VarDecls) {
	fmt.Fprintf(p, "var ")
	if d.Lparen == nil {
		// single declare
		for _, decl := range d.Decls {
			printVarDecl(p, decl)
		}
	} else {
		fmt.Fprintf(p, "(\n")
		p.Tab()
		for _, decl := range d.Decls {
			printVarDecl(p, decl)
			fmt.Println(p)
		}
		p.ShiftTab()
		fmt.Fprintf(p, ")")
	}
}

func printConstDecl(p *fmt8.Printer, d *ast.ConstDecl) {
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

func printConstDecls(p *fmt8.Printer, d *ast.ConstDecls) {
	fmt.Fprintf(p, "const ")
	if d.Lparen == nil {
		// single declare
		for _, decl := range d.Decls {
			printConstDecl(p, decl)
		}
	} else {
		fmt.Fprintf(p, "(\n")
		p.Tab()
		for _, decl := range d.Decls {
			printConstDecl(p, decl)
			fmt.Println(p)
		}
		p.ShiftTab()
		fmt.Fprintf(p, ")")
	}
}

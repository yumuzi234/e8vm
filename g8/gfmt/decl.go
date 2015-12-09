package gfmt

import (
	"fmt"

	"e8vm.io/e8vm/fmt8"
	"e8vm.io/e8vm/g8/ast"
)

func printIdents(p *fmt8.Printer, m *matcher, idents *ast.IdentList) {
	for i, id := range idents.Idents {
		if i > 0 {
			printExprs(p, m, idents.Commas[i-1], " ")
		}
		printToken(p, m, id)
	}
}

func printVarDecl(p *fmt8.Printer, m *matcher, d *ast.VarDecl) {
	printIdents(p, m, d.Idents)
	if d.Type != nil {
		printExprs(p, m, " ", d.Type)
	}

	if d.Eq != nil {
		printExprs(p, m, " ", d.Eq, " ", d.Exprs)
	}
	fmt.Fprintln(p)
}

func printVarDecls(p *fmt8.Printer, m *matcher, d *ast.VarDecls) {
	printExprs(p, m, d.Kw, " ")
	if d.Lparen == nil {
		// single declare
		for _, decl := range d.Decls {
			printVarDecl(p, m, decl)
		}
	} else {
		printToken(p, m, d.Lparen)
		fmt.Fprintln(p)
		p.Tab()
		for _, decl := range d.Decls {
			printVarDecl(p, m, decl)
		}
		p.ShiftTab()
		printToken(p, m, d.Rparen)
		fmt.Fprintln(p)
	}
}

func printConstDecl(p *fmt8.Printer, m *matcher, d *ast.ConstDecl) {
	printIdents(p, m, d.Idents)
	if d.Type != nil {
		printExprs(p, m, " ", d.Type)
	}

	if d.Eq != nil {
		printExprs(p, m, " ", d.Eq, " ", d.Exprs)
	}
	fmt.Fprintln(p)
}

func printConstDecls(p *fmt8.Printer, m *matcher, d *ast.ConstDecls) {
	printExprs(p, m, d.Kw, " ")
	if d.Lparen == nil {
		// single declare
		for _, decl := range d.Decls {
			printConstDecl(p, m, decl)
		}
	} else {
		printToken(p, m, d.Lparen)
		fmt.Fprintln(p)
		p.Tab()
		for _, decl := range d.Decls {
			printConstDecl(p, m, decl)
		}
		p.ShiftTab()
		printToken(p, m, d.Rparen)
		fmt.Fprintln(p)
	}
}

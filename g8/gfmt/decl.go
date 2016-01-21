package gfmt

import (
	"e8vm.io/e8vm/g8/ast"
)

func printIdents(f *formatter, idents *ast.IdentList) {
	for i, id := range idents.Idents {
		if i > 0 {
			printExprs(f, idents.Commas[i-1], " ")
		}
		f.printToken(id)
	}
}

func printVarDecl(f *formatter, d *ast.VarDecl) {
	printIdents(f, d.Idents)
	if d.Type != nil {
		printExprs(f, " ", d.Type)
	}
	if d.Eq != nil {
		printExprs(f, " ", d.Eq, " ", d.Exprs)
	}
}

func printVarDecls(f *formatter, d *ast.VarDecls) {
	printExprs(f, d.Kw, " ")
	if d.Lparen == nil {
		// single declare
		printVarDecl(f, d.Decls[0])
	} else {
		f.printToken(d.Lparen)
		f.printEndl()
		f.Tab()
		for i, decl := range d.Decls {
			printVarDecl(f, decl)
			f.printEndlPlus(i < len(d.Decls)-1, 0)
		}
		f.ShiftTab()
		f.printToken(d.Rparen)
	}
}

func printConstDecl(f *formatter, d *ast.ConstDecl) {
	printIdents(f, d.Idents)
	if d.Type != nil {
		printExprs(f, " ", d.Type)
	}
	if d.Eq != nil {
		printExprs(f, " ", d.Eq, " ", d.Exprs)
	}
}

func printConstDecls(f *formatter, d *ast.ConstDecls) {
	printExprs(f, d.Kw, " ")
	if d.Lparen == nil {
		// single declare
		printConstDecl(f, d.Decls[0])
	} else {
		f.printToken(d.Lparen)
		f.printEndl()
		f.Tab()
		for i, decl := range d.Decls {
			printConstDecl(f, decl)
			f.printEndlPlus(i < len(d.Decls)-1, 0)
		}
		f.ShiftTab()
		f.printToken(d.Rparen)
	}
}

func printImportDecl(f *formatter, d *ast.ImportDecl) {
	if d.As != nil {
		printExprs(f, d.As, " ")
	}
	printExprs(f, d.Path)
}

func printImportDecls(f *formatter, d *ast.ImportDecls) {
	printExprs(f, d.Kw, " ")
	f.printToken(d.Lparen)
	f.printEndl()
	f.Tab()
	// TODO: sort imports in groups
	for i, decl := range d.Decls {
		printImportDecl(f, decl)
		f.printEndlPlus(i < len(d.Decls)-1, 0)
	}
	f.ShiftTab()
	f.printToken(d.Rparen)
}

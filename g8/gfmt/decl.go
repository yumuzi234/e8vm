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
	f.printEndl()
}

func printVarDecls(f *formatter, d *ast.VarDecls) {
	printExprs(f, d.Kw, " ")
	if d.Lparen == nil {
		// single declare
		for _, decl := range d.Decls {
			printVarDecl(f, decl)
		}
	} else {
		f.printToken(d.Lparen)
		f.printEndl()
		f.Tab()
		for _, decl := range d.Decls {
			printVarDecl(f, decl)
		}
		f.ShiftTab()
		f.printToken(d.Rparen)
		f.printEndl()
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
	f.printEndl()
}

func printConstDecls(f *formatter, d *ast.ConstDecls) {
	printExprs(f, d.Kw, " ")
	if d.Lparen == nil {
		// single declare
		for _, decl := range d.Decls {
			printConstDecl(f, decl)
		}
	} else {
		f.printToken(d.Lparen)
		f.printEndl()
		f.Tab()
		for _, decl := range d.Decls {
			printConstDecl(f, decl)
		}
		f.ShiftTab()
		f.printToken(d.Rparen)
		f.printEndl()
	}
}

func printImportDecl(f *formatter, d *ast.ImportDecl) {
	if d.As != nil {
		printExprs(f, d.As, " ")
	}
	printExprs(f, d.Path)
	f.printEndl()
}

func printImportDecls(f *formatter, d *ast.ImportDecls) {
	printExprs(f, d.Kw, " ")
	f.printToken(d.Lparen)
	f.printEndl()
	f.Tab()
	// TODO: sort imports in groups
	for _, decl := range d.Decls {
		printImportDecl(f, decl)
	}
	f.ShiftTab()
	f.printToken(d.Rparen)
	f.printEndl()
}

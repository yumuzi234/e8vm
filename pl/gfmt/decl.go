package gfmt

import (
	"shanhu.io/smlvm/pl/ast"
)

func printIdents(f *formatter, idents *ast.IdentList) {
	for i, id := range idents.Idents {
		if i > 0 {
			f.printExprs(idents.Commas[i-1], " ")
		}
		f.printToken(id)
	}
}

func printVarDecl(f *formatter, d *ast.VarDecl) {
	printIdents(f, d.Idents)
	if d.Type != nil {
		f.printExprs(" ", d.Type)
	}
	if d.Eq != nil {
		f.printExprs(" ", d.Eq, " ", d.Exprs)
	}
}

func printVarDecls(f *formatter, d *ast.VarDecls) {
	f.printExprs(d.Kw, " ")
	if d.Lparen == nil {
		// single declare
		printVarDecl(f, d.Decls[0])
	} else {
		f.printToken(d.Lparen)
		f.printEndl()
		f.Tab()
		for i, decl := range d.Decls {
			if i > 0 {
				f.printGap()
			}
			printVarDecl(f, decl)
		}
		f.printEndl()
		f.ShiftTab()
		f.printToken(d.Rparen)
	}
}

func printConstDecl(f *formatter, d *ast.ConstDecl) {
	printIdents(f, d.Idents)
	if d.Type != nil {
		f.printExprs(" ", d.Type)
	}
	if d.Eq != nil {
		f.printExprs(" ", d.Eq, " ", d.Exprs)
	}
}

func printConstDecls(f *formatter, d *ast.ConstDecls) {
	f.printExprs(d.Kw, " ")
	if d.Lparen == nil {
		// single declare
		printConstDecl(f, d.Decls[0])
	} else {
		f.printToken(d.Lparen)
		f.printEndl()
		f.Tab()
		for i, decl := range d.Decls {
			if i > 0 {
				f.printGap()
			}
			printConstDecl(f, decl)
		}
		f.printEndl()
		f.ShiftTab()
		f.printToken(d.Rparen)
	}
}

func printImportDecl(f *formatter, d *ast.ImportDecl) {
	if d.As != nil {
		f.printExprs(d.As, " ")
	}
	f.printExprs(d.Path)
}

func printImportDecls(f *formatter, d *ast.ImportDecls) {
	f.printExprs(d.Kw, " ")
	f.printToken(d.Lparen)
	f.printEndl()
	f.Tab()
	// TODO(yumuzi): sort imports in groups
	for i, decl := range d.Decls {
		if i > 0 {
			f.printGap()
		}
		printImportDecl(f, decl)
	}
	f.printEndl()
	f.ShiftTab()
	f.printToken(d.Rparen)
}

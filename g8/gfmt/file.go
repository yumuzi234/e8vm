package gfmt

import (
	"e8vm.io/e8vm/g8/ast"
)

func printTopDecl(f *formatter, d ast.Decl) {
	switch d := d.(type) {
	case *ast.Func:
		printFunc(f, d)
	case *ast.Struct:
		printStruct(f, d)
	case *ast.VarDecls:
		printVarDecls(f, d)
	case *ast.ConstDecls:
		printConstDecls(f, d)
	default:
		f.errorf(nil, "invalid top-level declaration type: %T", d)
	}
}

func printFile(f *formatter, file *ast.File) {
	if file.Imports != nil {
		printImportDecls(f, file.Imports)
		f.printEndlPlus(true, true)
	}

	for i, decl := range file.Decls {
		printTopDecl(f, decl)
		f.printEndlPlus(i < len(file.Decls)-1, true)
	}
	f.finish()
}

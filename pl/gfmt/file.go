package gfmt

import (
	"shanhu.io/smlvm/pl/ast"
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
	case *ast.Interface:
		printInterface(f, d)
	default:
		f.errorf(nil, "invalid top-level declaration type: %T", d)
	}
}

func printFile(f *formatter, file *ast.File) {
	if file.Imports != nil {
		printImportDecls(f, file.Imports)
		f.printEndlPlus(len(file.Decls) > 0, true)
	}

	// empty line between each topDecl?
	for i, decl := range file.Decls {
		printTopDecl(f, decl)
		f.printEndlPlus(i < len(file.Decls)-1, true)
	}
	f.finish()
}

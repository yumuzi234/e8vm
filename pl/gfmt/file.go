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
		f.printEndl()
		if len(file.Decls) != 0 {
			f.printEndl()
		}
	}
	if len(file.Decls) == 0 {
		f.finish()
		return
	}
	for i, decl := range file.Decls {
		if i != 0 {
			f.printEndl()
		}
		printTopDecl(f, decl)
		f.printEndl()
	}
	f.finish()
}

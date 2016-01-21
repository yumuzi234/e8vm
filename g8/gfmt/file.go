package gfmt

import (
	"fmt"
	"io"

	"e8vm.io/e8vm/g8/ast"
	"e8vm.io/e8vm/lex8"
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
		panic(fmt.Errorf("invalid top-level declaration type: %T", d))
	}
}

func printFile(f *formatter, file *ast.File) {
	if file.Imports != nil {
		printImportDecls(f, file.Imports)
		f.printEndlPlus(true, 1)
	}

	for i, decl := range file.Decls {
		printTopDecl(f, decl)
		f.printEndlPlus(i < len(file.Decls)-1, 1)
	}
	f.finish()
}

// FprintFile prints a file
func FprintFile(out io.Writer, file *ast.File, rec *lex8.Recorder) {
	f := newFormatter(out, rec.Tokens())
	printFile(f, file)
}

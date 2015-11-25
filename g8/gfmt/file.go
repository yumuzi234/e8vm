package gfmt

import (
	"fmt"
	"io"

	"e8vm.io/e8vm/fmt8"
	"e8vm.io/e8vm/g8/ast"
)

func printTopDecl(p *fmt8.Printer, d ast.Decl) {
	switch d := d.(type) {
	case *ast.Func:
		printFunc(p, d)
	case *ast.Struct:
		printStruct(p, d)
	case *ast.VarDecls:
		printVarDecls(p, d)
	case *ast.ConstDecls:
		printConstDecls(p, d)
	default:
		panic(fmt.Errorf("invalid top-level declaration type: %T", d))
	}
}

func printFile(p *fmt8.Printer, f *ast.File) {
	for i, decl := range f.Decls {
		printTopDecl(p, decl)
		if i < len(f.Decls)-1 {
			fmt.Fprintln(p)
		}
	}
}

// FprintFile prints a list of file
func FprintFile(out io.Writer, f *ast.File) {
	p := fmt8.NewPrinter(out)
	printFile(p, f)
}

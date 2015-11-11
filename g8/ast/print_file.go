package ast

import (
	"fmt"
	"io"

	"e8vm.io/e8vm/fmt8"
)

func printTopDecl(p *fmt8.Printer, d Decl) {
	switch d := d.(type) {
	case *Func:
		printFunc(p, d)
	case *Struct:
		printStruct(p, d)
	case *VarDecls:
		printVarDecls(p, d)
	case *ConstDecls:
		printConstDecls(p, d)
	default:
		panic(fmt.Errorf("invalid top-level declaration type: %T", d))
	}
}

func printFile(p *fmt8.Printer, f *File) {
	for i, decl := range f.Decls {
		printTopDecl(p, decl)
		if i < len(f.Decls)-1 {
			fmt.Fprintln(p)
		}
	}
}

// FprintFile prints a list of file
func FprintFile(out io.Writer, f *File) {
	p := fmt8.NewPrinter(out)
	printFile(p, f)
}

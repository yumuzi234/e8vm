package gfmt

import (
	"fmt"
	"io"

	"e8vm.io/e8vm/fmt8"
	"e8vm.io/e8vm/g8/ast"
	"e8vm.io/e8vm/lex8"
)

func printTopDecl(p *fmt8.Printer, m *matcher, d ast.Decl) {
	switch d := d.(type) {
	case *ast.Func:
		printFunc(p, m, d)
	case *ast.Struct:
		printStruct(p, m, d)
	case *ast.VarDecls:
		printVarDecls(p, m, d)
	case *ast.ConstDecls:
		printConstDecls(p, m, d)
	default:
		panic(fmt.Errorf("invalid top-level declaration type: %T", d))
	}
}

func printFile(p *fmt8.Printer, m *matcher, f *ast.File) {
	for i, decl := range f.Decls {
		printTopDecl(p, m, decl)
		if i < len(f.Decls)-1 {
			fmt.Fprintln(p)
		}
	}
	m.finish()
}

// FprintFile prints a file
func FprintFile(out io.Writer, f *ast.File, rec *lex8.Recorder) {
	p := fmt8.NewPrinter(out)
	m := newMatcher(rec.Tokens())
	printFile(p, m, f)
}

package asm

import (
	"e8vm.io/e8vm/asm/ast"
	"e8vm.io/e8vm/lexing"
)

type file struct {
	*ast.File

	funcs   []*funcDecl
	vars    []*varDecl
	imports *importDecl
}

func resolveFile(log lexing.Logger, f *ast.File) *file {
	ret := new(file)
	ret.File = f

	if f.Imports != nil {
		ret.imports = resolveImportDecl(log, f.Imports)
	}

	for _, d := range f.Decls {
		if d, ok := d.(*ast.Func); ok {
			ret.funcs = append(ret.funcs, resolveFunc(log, d))
		}

		if d, ok := d.(*ast.Var); ok {
			ret.vars = append(ret.vars, resolveVar(log, d))
		}
	}

	return ret
}

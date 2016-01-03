package sempass

import (
	"e8vm.io/e8vm/asm8"
	"e8vm.io/e8vm/g8/ast"
	"e8vm.io/e8vm/g8/tast"
	"e8vm.io/e8vm/g8/types"
	"e8vm.io/e8vm/sym8"
)

func declareFuncSym(b *Builder, f *ast.Func, t types.T) *sym8.Symbol {
	name := f.Name.Lit
	s := sym8.Make(b.path, name, tast.SymFunc, nil, t, f.Name.Pos)
	conflict := b.scope.Declare(s)
	if conflict != nil {
		b.Errorf(f.Name.Pos, "%q already defined as a %s",
			name, tast.SymStr(conflict.Type),
		)
		b.Errorf(conflict.Pos, "previously defined here")
		return nil
	}
	return s
}

func buildFuncAlias(b *Builder, f *ast.Func, t *types.Func) *sym8.Symbol {
	alias := f.Alias
	pkg := buildPkgRef(b, alias.Pkg)
	if pkg == nil {
		return nil
	}
	if pkg.Lang != "asm8" {
		b.Errorf(alias.Pkg.Pos, "can only alias functions in asm packages")
		return nil
	}

	sym := findPackageSym(b, alias.Name, pkg)
	if sym == nil {
		return nil
	}

	if sym.Type != asm8.SymFunc {
		b.Errorf(alias.Name.Pos, "the symbol is not a function")
		return nil
	}

	ret := declareFuncSym(b, f, t)
	return ret
}

func declareFunc(b *Builder, f *ast.Func) *sym8.Symbol {
	t := buildFuncType(b, nil, f.FuncSig)
	if t == nil {
		return nil
	}

	if f.Alias != nil {
		return buildFuncAlias(b, f, t)
	}

	panic("todo")
}

func buildFuncs(b *Builder, funcs []*ast.Func) []*sym8.Symbol {

	panic("todo")
}

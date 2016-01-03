package sempass

import (
	"e8vm.io/e8vm/asm8"
	"e8vm.io/e8vm/g8/ast"
	"e8vm.io/e8vm/g8/tast"
	"e8vm.io/e8vm/g8/types"
	"e8vm.io/e8vm/sym8"
)

type pkgFunc struct {
	sym *sym8.Symbol
	ast *ast.Func
}

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

func buildFuncAlias(b *Builder, f *ast.Func) *tast.FuncAlias {
	t := buildFuncType(b, nil, f.FuncSig)
	if t == nil {
		return nil
	}

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
	if ret == nil {
		return nil
	}

	return &tast.FuncAlias{Sym: ret, Of: sym}
}

func declareFunc(b *Builder, f *ast.Func) *pkgFunc {
	if f.Alias != nil {
		panic("bug")
	}

	t := buildFuncType(b, nil, f.FuncSig)
	if t == nil {
		return nil
	}

	s := declareFuncSym(b, f, t)
	if s == nil {
		return nil
	}

	return &pkgFunc{sym: s, ast: f}
}

func declareFuncs(b *Builder, funcs []*ast.Func) (
	[]*pkgFunc, []*tast.FuncAlias,
) {
	var ret []*pkgFunc
	var aliases []*tast.FuncAlias

	for _, f := range funcs {
		if f.Alias != nil {
			a := buildFuncAlias(b, f)
			if a != nil {
				aliases = append(aliases, a)
			}
			continue
		}

		r := declareFunc(b, f)
		if r != nil {
			ret = append(ret, r)
		}
	}

	return ret, aliases
}

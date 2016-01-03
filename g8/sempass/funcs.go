package sempass

import (
	"e8vm.io/e8vm/g8/ast"
	"e8vm.io/e8vm/g8/tast"
)

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

func buildFuncs(b *Builder, funcs []*pkgFunc) []*tast.Func {
	b.this = nil
	b.thisType = nil

	ret := make([]*tast.Func, 0, len(funcs))
	for _, f := range funcs {
		res := buildFunc(b, f)
		if res != nil {
			ret = append(ret, res)
		}
	}

	return ret
}

func declareMethods(
	b *Builder, methods []*ast.Func, pkgStructs []*pkgStruct,
) []*pkgFunc {
	m := make(map[string]*pkgStruct)
	for _, ps := range pkgStructs {
		m[ps.name.Lit] = ps
	}

	var ret []*pkgFunc

	// inlined ones
	for _, ps := range pkgStructs {
		for _, f := range ps.ast.Methods {
			pf := declareMethod(b, ps, f)
			if pf != nil {
				ret = append(ret, pf)
			}
		}
	}

	// go-like ones
	for _, f := range methods {
		recv := f.Recv.StructName
		ps := m[recv.Lit]
		if ps != nil {
			b.Errorf(recv.Pos, "struct %s not defined", recv.Lit)
			continue
		}

		pf := declareMethod(b, ps, f)
		if pf != nil {
			ret = append(ret, pf)
		}
	}

	return ret
}

func buildMethods(b *Builder, funcs []*pkgFunc) []*tast.Func {
	var ret []*tast.Func
	for _, f := range funcs {
		r := buildMethod(b, f)
		if r != nil {
			ret = append(ret, r)
		}
	}

	return ret
}

package sempass

import (
	"e8vm.io/e8vm/builds"
	"e8vm.io/e8vm/glang/ast"
	"e8vm.io/e8vm/glang/tast"
	"e8vm.io/e8vm/glang/types"
	"e8vm.io/e8vm/syms"
)

func buildImports(
	b *builder, f *ast.File, imps map[string]*builds.Package,
) []*syms.Symbol {
	if f.Imports == nil {
		return nil
	}

	var ret []*syms.Symbol
	for _, d := range f.Imports.Decls {
		_, as, e := ast.ImportPathAs(d)
		if e != nil {
			b.Errorf(d.Path.Pos, "invalid import path")
			continue
		}

		p := imps[as]
		if p == nil {
			b.Errorf(d.Path.Pos, "package %s missing", as)
			continue
		}

		pos := ast.ImportPos(d)
		if !(p.Lang == "asm8" || p.Lang == "g8") {
			b.Errorf(pos, "cannot import %q pacakge %q",
				p.Lang, as,
			)
			continue
		}

		t := &types.Pkg{as, p.Lang, p.Symbols}
		sym := syms.Make(b.path, as, tast.SymImport, nil, t, pos)
		conflict := b.scope.Declare(sym)
		if conflict != nil {
			b.Errorf(pos, "%s already imported", as)
			continue
		}

		ret = append(ret, sym)
	}
	return ret
}

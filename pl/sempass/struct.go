package sempass

import (
	"shanhu.io/smlvm/lexing"
	"shanhu.io/smlvm/pl/ast"
	"shanhu.io/smlvm/pl/tast"
	"shanhu.io/smlvm/pl/types"
	"shanhu.io/smlvm/syms"
)

type pkgStruct struct {
	name *lexing.Token
	ast  *ast.Struct    // the struct AST node
	sym  *syms.Symbol   // the symbol
	t    *types.Struct  // type
	pt   *types.Pointer // pointer type
	deps []string       // depending identifiers
}

func newPkgStruct(s *ast.Struct) *pkgStruct {
	deps := listStructDeps(s)
	t := types.NewStruct(s.Name.Lit)

	return &pkgStruct{
		name: s.Name,
		ast:  s,
		deps: deps,
		t:    t,
		pt:   types.NewPointer(t),
	}
}

func declareStruct(b *builder, s *ast.Struct) *pkgStruct {
	ret := newPkgStruct(s)
	name := ret.name.Lit
	pos := ret.name.Pos
	t := &types.Type{T: ret.t}
	sym := syms.Make(b.path, name, tast.SymStruct, nil, t, pos)
	conflict := b.scope.Declare(sym)
	if conflict != nil {
		b.CodeErrorf(pos, "pl.declConflict.struct",
			"%s already defined", name)
		b.CodeErrorf(conflict.Pos, "pl.declConflict.previousPos",
			"previously defined here as a %s", tast.SymStr(conflict.Type))
		return nil
	}

	ret.sym = sym
	return ret
}

func sortStructs(b *builder, m map[string]*pkgStruct) []*pkgStruct {
	s := newTopoSorter("struct", "pl.circDep.struct")
	for name, ps := range m {
		s.addNode(name, ps.name, ps.deps)
	}

	order := s.sort(b)
	ret := make([]*pkgStruct, 0, len(order))
	for _, name := range order {
		ret = append(ret, m[name])
	}
	return ret
}

func buildFields(b *builder, ps *pkgStruct) {
	s := ps.ast
	t := ps.t

	for _, f := range s.Fields {
		ft := b.buildType(f.Type)
		if ft == nil {
			continue
		}

		for _, id := range f.Idents.Idents {
			name := id.Lit
			field := &types.Field{Name: name, T: ft}
			sym := syms.Make(b.path, name, tast.SymField, field, ft, id.Pos)
			conflict := t.Syms.Declare(sym)
			if conflict != nil {
				b.CodeErrorf(id.Pos, "pl.declConflict.field",
					"field %s already defined", id.Lit)
				b.CodeErrorf(conflict.Pos,
					"pl.declConflict.previousPos",
					"previously defined here")
				continue
			}

			t.AddField(field)
		}
	}
}

func declareStructs(b *builder, structs []*ast.Struct) []*pkgStruct {
	m := make(map[string]*pkgStruct)
	for _, s := range structs {
		ps := declareStruct(b, s)
		if ps != nil {
			m[ps.name.Lit] = ps
		}
	}
	ret := sortStructs(b, m)
	return ret
}

func buildStructs(b *builder, pkystructs []*pkgStruct) {
	for _, ps := range pkystructs {
		buildFields(b, ps)
	}
}

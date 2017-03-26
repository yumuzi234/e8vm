package sempass

import (
	"shanhu.io/smlvm/lexing"
	"shanhu.io/smlvm/pl/ast"
	"shanhu.io/smlvm/pl/tast"
	"shanhu.io/smlvm/pl/types"
	"shanhu.io/smlvm/syms"
)

type pkgInterface struct {
	name *lexing.Token
	ast  *ast.Interface   // the AST node
	sym  *syms.Symbol     // the symbol
	t    *types.Interface // type

	// this not needed for interface, right?
	// pt   *types.Pointer   // pointer type
}

func newPkgInterface(s *ast.Interface) *pkgInterface {
	t := types.NewInterface(s.Name.Lit)
	return &pkgInterface{
		name: s.Name,
		ast:  s,
		t:    t,
	}
}

func declareInterface(b *builder, i *ast.Interface) *pkgInterface {
	ret := newPkgInterface(i)
	name := ret.name.Lit
	pos := ret.name.Pos
	t := &types.Type{T: ret.t}
	sym := syms.Make(b.path, name, tast.SymInterface, nil, t, pos)
	conflict := b.scope.Declare(sym)
	if conflict != nil {
		b.CodeErrorf(pos, "pl.declConflict.interface",
			"%s already defined", name)
		b.CodeErrorf(conflict.Pos, "pl.declConflict.previousPos",
			"previously defined here as a %s", tast.SymStr(conflict.Type))
		return nil
	}

	ret.sym = sym
	return ret
}

func declareInterfaces(b *builder, is []*ast.Interface) []*pkgInterface {
	ret := make([]*pkgInterface, 0)
	for _, i := range is {
		pi := declareInterface(b, i)
		if pi != nil {
			ret = append(ret, pi)
		}
	}
	return ret
}

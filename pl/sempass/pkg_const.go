package sempass

import (
	"shanhu.io/smlvm/lexing"
	"shanhu.io/smlvm/pl/ast"
	"shanhu.io/smlvm/pl/tast"
	"shanhu.io/smlvm/pl/types"
	"shanhu.io/smlvm/syms"
	"shanhu.io/smlvm/toposort"
)

type pkgConst struct {
	sym  *syms.Symbol
	tok  *lexing.Token
	expr ast.Expr
	deps []string
}

func buildPkgConstDecl(b *builder, d *ast.ConstDecl) []*pkgConst {
	if d.Type != nil {
		b.Errorf(ast.ExprPos(d.Type), "typed const not implemented")
		return nil
	}

	nident := len(d.Idents.Idents)
	nexpr := len(d.Exprs.Exprs)
	if nident != nexpr {
		b.Errorf(d.Eq.Pos, "%d consts with %d expressions",
			nident, nexpr,
		)
		return nil
	}

	var ret []*pkgConst
	zero := types.NewNumber(0) // place holding the const
	for i, ident := range d.Idents.Idents {
		s := declareConst(b, ident, zero)
		if s == nil {
			return nil
		}

		expr := d.Exprs.Exprs[i]
		c := &pkgConst{
			sym:  s,
			tok:  ident,
			expr: expr,
			deps: symUse(expr),
		}
		ret = append(ret, c)
	}

	return ret
}

func sortPkgConsts(b *builder, consts []*pkgConst) []*pkgConst {
	m := make(map[string]*pkgConst)
	s := toposort.NewSorter("const")
	for _, c := range consts {
		name := c.sym.Name()
		m[name] = c
		s.AddNode(name, c.tok, c.deps)
	}

	order := s.Sort(b)
	var ret []*pkgConst
	for _, name := range order {
		ret = append(ret, m[name])
	}
	return ret
}

func buildPkgConst(b *builder, c *pkgConst) *syms.Symbol {
	right := b.buildConstExpr(c.expr)
	if right == nil {
		return nil
	}
	res, ok := right.(*tast.Const)
	if !ok {
		b.Errorf(ast.ExprPos(c.expr), "expect a const")
		return nil
	}

	t := res.Ref.Type()
	c.sym.ObjType = t

	return c.sym
}

func buildPkgConsts(b *builder, consts []*ast.ConstDecls) []*syms.Symbol {
	var res []*pkgConst
	for _, c := range consts {
		for _, d := range c.Decls {
			pkgConsts := buildPkgConstDecl(b, d)
			if pkgConsts != nil {
				res = append(res, pkgConsts...)
			}
		}
	}

	var ret []*syms.Symbol
	res = sortPkgConsts(b, res)
	for _, c := range res {
		s := buildPkgConst(b, c)
		if s != nil {
			ret = append(ret, s)
		}
	}
	return ret
}

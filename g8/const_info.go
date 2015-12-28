package g8

import (
	"e8vm.io/e8vm/g8/ast"
	"e8vm.io/e8vm/lex8"
	"e8vm.io/e8vm/toposort"
)

type constInfo struct {
	name *lex8.Token
	typ  ast.Expr
	expr ast.Expr

	deps []string
}

func newConstInfo(name *lex8.Token, typ, expr ast.Expr) *constInfo {
	return &constInfo{
		name: name, typ: typ, expr: expr,
		deps: symUseExpr(expr),
	}
}

func sortConsts(b *builder, m map[string]*constInfo) []*constInfo {
	s := toposort.NewSorter("const")
	for name, info := range m {
		s.AddNode(name, info.name, info.deps)
	}

	order := s.Sort(b)
	var ret []*constInfo
	for _, name := range order {
		ret = append(ret, m[name])
	}

	return ret
}

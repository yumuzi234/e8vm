package sempass

import (
	"sort"

	"e8vm.io/e8vm/glang/ast"
	"e8vm.io/e8vm/glang/parse"
)

type structDeps struct {
	deps map[string]struct{}
}

func newStructDeps() *structDeps {
	return &structDeps{
		deps: make(map[string]struct{}),
	}
}

func (d *structDeps) add(t ast.Expr) {
	switch t := t.(type) {
	case *ast.Operand:
		if t.Token.Type == parse.Ident {
			d.deps[t.Token.Lit] = struct{}{}
		}
	case *ast.ParenExpr:
		d.add(t.Expr)
	case *ast.ArrayTypeExpr:
		if t.Len != nil {
			d.add(t.Type)
		}
	}
}

func (d *structDeps) list() []string {
	ret := make([]string, 0, len(d.deps))
	for k := range d.deps {
		ret = append(ret, k)
	}
	sort.Strings(ret)
	return ret
}

func listStructDeps(s *ast.Struct) []string {
	deps := newStructDeps()
	for _, f := range s.Fields {
		deps.add(f.Type)
	}
	return deps.list()
}

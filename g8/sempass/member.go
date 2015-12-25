package sempass

import (
	"e8vm.io/e8vm/g8/ast"
	"e8vm.io/e8vm/g8/tast"
	"e8vm.io/e8vm/g8/types"
	"e8vm.io/e8vm/lex8"
	"e8vm.io/e8vm/sym8"
)

func findPackageSym(
	b *Builder, sub *lex8.Token, pkg *types.Pkg,
) *sym8.Symbol {
	sym := pkg.Syms.Query(sub.Lit)
	if sym == nil {
		b.Errorf(sub.Pos, "%s has no symbol named %s",
			pkg, sub.Lit,
		)
		return nil
	}
	name := sym.Name()
	if !sym8.IsPublic(name) && sym.Pkg() != b.path {
		b.Errorf(sub.Pos, "symbol %s is not public", name)
		return nil
	}

	return sym
}

func buildConstMember(b *Builder, m *ast.MemberExpr) tast.Expr {
	obj := b.BuildConstExpr(m.Expr)
	if obj == nil {
		return nil
	}

	ref := tast.ExprRef(obj)
	if ref.List != nil {
		b.Errorf(m.Dot.Pos, "expression list does not have any member")
		return nil
	}

	if pkg, ok := ref.T.(*types.Pkg); ok {
		s := findPackageSym(b, m.Sub, pkg)
		if s == nil {
			return nil
		}
		if s.Type != tast.SymConst {
			b.Errorf(m.Sub.Pos, "%s.%s is not a const", pkg, m.Sub.Lit)
			return nil
		}
		return &tast.Const{tast.NewRef(s.ObjType.(types.T))}
	}

	b.Errorf(m.Dot.Pos, "expect const expression")
	return nil
}

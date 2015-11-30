package g8

import (
	"e8vm.io/e8vm/g8/ast"
	"e8vm.io/e8vm/g8/parse"
	"e8vm.io/e8vm/lex8"
)

type symUses struct {
	uses map[string]*lex8.Token
}

func newSymUses() *symUses {
	return &symUses{
		uses: make(map[string]*lex8.Token),
	}
}

func (u *symUses) add(tok *lex8.Token) {
	name := tok.Lit
	if _, found := u.uses[name]; !found {
		u.uses[name] = tok
	}
}

func (u *symUses) symUseExpr(expr ast.Expr) {
	if expr == nil {
		return
	}

	switch expr := expr.(type) {
	case *ast.Operand:
		if expr.Token.Type == parse.Ident {
			u.add(expr.Token)
		}
	case *ast.MemberExpr:
		u.symUseExpr(expr.Expr)
	case *ast.OpExpr:
		u.symUseExpr(expr.A)
		u.symUseExpr(expr.B)
	case *ast.ParenExpr:
		u.symUseExpr(expr.Expr)
	case *ast.StarExpr:
		u.symUseExpr(expr.Expr)
	case *ast.CallExpr:
		u.symUseExpr(expr.Func)
		for _, arg := range expr.Args.Exprs {
			u.symUseExpr(arg)
		}
	case *ast.IndexExpr:
		u.symUseExpr(expr.Array)
		u.symUseExpr(expr.Index)
		u.symUseExpr(expr.IndexEnd)
	case *ast.ArrayTypeExpr:
		u.symUseExpr(expr.Len)
		u.symUseExpr(expr.Type)
	case *ast.FuncTypeExpr:
		sig := expr.FuncSig
		for _, arg := range sig.Args.Paras {
			u.symUseExpr(arg.Type)
		}
		for _, ret := range sig.Rets.Paras {
			u.symUseExpr(ret.Type)
		}
		u.symUseExpr(sig.RetType)
	}
}

func symUseExpr(expr ast.Expr) []string {
	u := newSymUses()
	u.symUseExpr(expr)

	var ret []string
	for name := range u.uses {
		ret = append(ret, name)
	}
	return ret
}

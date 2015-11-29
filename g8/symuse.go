package g8

import (
	"e8vm.io/e8vm/g8/ast"
	"e8vm.io/e8vm/g8/parse"
	"e8vm.io/e8vm/lex8"
)

type symUses struct {
	toks []*lex8.Token
}

func newSymUses() *symUses {
	return &symUses{}
}

func (u *symUses) add(tok *lex8.Token) {
	u.toks = append(u.toks, tok)
}

func symUseExpr(u *symUses, expr ast.Expr) {
	if expr == nil {
		return
	}

	switch expr := expr.(type) {
	case *ast.Operand:
		if expr.Token.Type == parse.Ident {
			u.add(expr.Token)
		}
	case *ast.MemberExpr:
		symUseExpr(u, expr.Expr)
	case *ast.OpExpr:
		symUseExpr(u, expr.A)
		symUseExpr(u, expr.B)
	case *ast.ParenExpr:
		symUseExpr(u, expr.Expr)
	case *ast.StarExpr:
		symUseExpr(u, expr.Expr)
	case *ast.CallExpr:
		symUseExpr(u, expr.Func)
		for _, arg := range expr.Args.Exprs {
			symUseExpr(u, arg)
		}
	case *ast.IndexExpr:
		symUseExpr(u, expr.Array)
		symUseExpr(u, expr.Index)
		symUseExpr(u, expr.IndexEnd)
	case *ast.ArrayTypeExpr:
		symUseExpr(u, expr.Len)
		symUseExpr(u, expr.Type)
	case *ast.FuncTypeExpr:
		sig := expr.FuncSig
		for _, arg := range sig.Args.Paras {
			symUseExpr(u, arg.Type)
		}
		for _, ret := range sig.Rets.Paras {
			symUseExpr(u, ret.Type)
		}
		symUseExpr(u, sig.RetType)
	}
}

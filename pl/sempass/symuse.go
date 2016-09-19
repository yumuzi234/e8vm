package sempass

import (
	"shanhu.io/smlvm/lexing"
	"shanhu.io/smlvm/pl/ast"
	"shanhu.io/smlvm/pl/parse"
)

type symUses struct {
	uses map[string]*lexing.Token
}

func newSymUses() *symUses {
	return &symUses{
		uses: make(map[string]*lexing.Token),
	}
}

func (u *symUses) add(tok *lexing.Token) {
	name := tok.Lit
	if _, found := u.uses[name]; !found {
		u.uses[name] = tok
	}
}

func (u *symUses) symUse(expr ast.Expr) {
	if expr == nil {
		return
	}

	switch expr := expr.(type) {
	case *ast.Operand:
		if expr.Token.Type == parse.Ident {
			u.add(expr.Token)
		}
	case *ast.MemberExpr:
		u.symUse(expr.Expr)
	case *ast.OpExpr:
		u.symUse(expr.A)
		u.symUse(expr.B)
	case *ast.ParenExpr:
		u.symUse(expr.Expr)
	case *ast.StarExpr:
		u.symUse(expr.Expr)
	case *ast.CallExpr:
		u.symUse(expr.Func)
		for _, arg := range expr.Args.Exprs {
			u.symUse(arg)
		}
	case *ast.IndexExpr:
		u.symUse(expr.Array)
		u.symUse(expr.Index)
		u.symUse(expr.IndexEnd)
	case *ast.ArrayTypeExpr:
		u.symUse(expr.Len)
		u.symUse(expr.Type)
	case *ast.FuncTypeExpr:
		sig := expr.FuncSig
		for _, arg := range sig.Args.Paras {
			u.symUse(arg.Type)
		}
		for _, ret := range sig.Rets.Paras {
			u.symUse(ret.Type)
		}
		u.symUse(sig.RetType)
	}
}

func symUse(expr ast.Expr) []string {
	u := newSymUses()
	u.symUse(expr)

	var ret []string
	for name := range u.uses {
		ret = append(ret, name)
	}
	return ret
}

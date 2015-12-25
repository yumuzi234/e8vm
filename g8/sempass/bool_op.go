package sempass

import (
	"e8vm.io/e8vm/g8/tast"
	"e8vm.io/e8vm/lex8"
)

func unaryOpBool(b *Builder, opTok *lex8.Token, B tast.Expr) tast.Expr {
	op := opTok.Lit
	if op == "!" {
		t := tast.ExprRef(B).T
		return &tast.OpExpr{nil, opTok, B, tast.NewRef(t)}
	}

	b.Errorf(opTok.Pos, "invalid operation: %q on boolean", op)
	return nil
}

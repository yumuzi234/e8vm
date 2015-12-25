package sempass

import (
	"e8vm.io/e8vm/g8/tast"
	"e8vm.io/e8vm/g8/types"
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

func binaryOpBool(b *Builder, opTok *lex8.Token, A, B tast.Expr) tast.Expr {
	op := opTok.Lit
	switch op {
	case "==", "!=":
		r := tast.NewRef(types.Bool)
		return &tast.OpExpr{A, opTok, B, r}
	}

	b.Errorf(opTok.Pos, "%q on bools", op)
	return nil
}

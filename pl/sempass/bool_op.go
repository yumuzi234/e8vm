package sempass

import (
	"shanhu.io/smlvm/lexing"
	"shanhu.io/smlvm/pl/tast"
	"shanhu.io/smlvm/pl/types"
)

func unaryOpBool(b *builder, opTok *lexing.Token, B tast.Expr) tast.Expr {
	op := opTok.Lit
	if op == "!" {
		t := B.R().T
		return &tast.OpExpr{Op: opTok, B: B, Ref: tast.NewRef(t)}
	}

	b.Errorf(opTok.Pos, "invalid operation: %q on boolean", op)
	return nil
}

func binaryOpBool(b *builder, opTok *lexing.Token, A, B tast.Expr) tast.Expr {
	op := opTok.Lit
	switch op {
	case "==", "!=", "&&", "||":
		r := tast.NewRef(types.Bool)
		return &tast.OpExpr{A: A, Op: opTok, B: B, Ref: r}
	}

	b.Errorf(opTok.Pos, "%q on bools", op)
	return nil
}

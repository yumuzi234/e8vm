package sempass

import (
	"shanhu.io/smlvm/lexing"
	"shanhu.io/smlvm/pl/tast"
	"shanhu.io/smlvm/pl/types"
)

func unaryOpInt(b *builder, opTok *lexing.Token, B tast.Expr) tast.Expr {
	op := opTok.Lit
	switch op {
	case "+":
		return B
	case "-", "^":
		t := B.R().T
		return &tast.OpExpr{Op: opTok, B: B, Ref: tast.NewRef(t)}
	}

	b.CodeErrorf(opTok.Pos, "pl.invalidOp",
		"invalid operation: %q on %s", op, B)
	return nil
}

func binaryOpInt(
	b *builder, opTok *lexing.Token, A, B tast.Expr, t types.T,
) tast.Expr {
	op := opTok.Lit
	switch op {
	case "+", "-", "*", "&", "|", "^", "%", "/":
		r := tast.NewRef(t)
		return &tast.OpExpr{A: A, Op: opTok, B: B, Ref: r}
	case "==", "!=", ">", "<", ">=", "<=":
		r := tast.NewRef(types.Bool)
		return &tast.OpExpr{A: A, Op: opTok, B: B, Ref: r}
	}

	b.Errorf(opTok.Pos, "%q on ints", op)
	return nil
}

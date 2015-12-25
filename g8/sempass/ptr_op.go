package sempass

import (
	"e8vm.io/e8vm/g8/tast"
	"e8vm.io/e8vm/g8/types"
	"e8vm.io/e8vm/lex8"
)

func refAddress(b *Builder, opTok *lex8.Token, B tast.Expr) tast.Expr {
	op := opTok.Lit
	opPos := opTok.Pos

	bref := tast.ExprRef(B)
	if types.IsType(bref.T) {
		b.Errorf(opPos, "%q on %s", op, bref.T)
		return nil
	} else if bref.List != nil {
		b.Errorf(opPos, "%q on expression list", op)
		return nil
	} else if !bref.Addressable {
		b.Errorf(opPos, "reading address of non-addressable")
		return nil
	}

	r := tast.NewRef(&types.Pointer{bref.T})
	return &tast.OpExpr{nil, opTok, B, r}
}

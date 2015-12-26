package sempass

import (
	"e8vm.io/e8vm/g8/tast"
	"e8vm.io/e8vm/g8/types"
	"e8vm.io/e8vm/lex8"
)

func refAddress(b *Builder, opTok *lex8.Token, B tast.Expr) tast.Expr {
	op := opTok.Lit
	opPos := opTok.Pos

	bref := B.R()
	if types.IsType(bref.T) || !bref.IsSingle() {
		b.Errorf(opPos, "%q on %s", op, bref)
		return nil
	} else if !bref.Addressable {
		b.Errorf(opPos, "reading address of non-addressable")
		return nil
	}

	r := tast.NewRef(&types.Pointer{bref.T})
	return &tast.OpExpr{nil, opTok, B, r}
}

func binaryOpNil(b *Builder, opTok *lex8.Token, A, B tast.Expr) tast.Expr {
	op := opTok.Lit
	switch op {
	case "==", "!=":
		return &tast.OpExpr{A, opTok, B, tast.NewRef(types.Bool)}
	}

	b.Errorf(opTok.Pos, "%q on nils", op)
	return nil
}

func binaryOpPtr(b *Builder, opTok *lex8.Token, A, B tast.Expr) tast.Expr {
	op := opTok.Lit
	atyp := A.R().T
	btyp := B.R().T

	switch op {
	case "==", "!=":
		if types.IsNil(atyp) {
			A = &tast.Cast{A, tast.NewRef(btyp)}
		} else if types.IsNil(btyp) {
			B = &tast.Cast{B, tast.NewRef(atyp)}
		}

		return &tast.OpExpr{A, opTok, B, tast.NewRef(types.Bool)}
	}

	b.Errorf(opTok.Pos, "%q on pointers", op)
	return nil
}

func binaryOpSlice(b *Builder, opTok *lex8.Token, A, B tast.Expr) tast.Expr {
	op := opTok.Lit
	switch op {
	case "==", "!=":
		return &tast.OpExpr{A, opTok, B, tast.NewRef(types.Bool)}
	}
	b.Errorf(opTok.Pos, "%q on slices", op)
	return nil
}

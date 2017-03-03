package sempass

import (
	"shanhu.io/smlvm/lexing"
	"shanhu.io/smlvm/pl/tast"
	"shanhu.io/smlvm/pl/types"
)

func unaryOpConst(b *builder, opTok *lexing.Token, B tast.Expr) tast.Expr {
	op := opTok.Lit
	bref := B.R()
	if !bref.IsSingle() {
		b.CodeErrorf(opTok.Pos, "pl.invalidExprStmt",
			"invalid expression, not single: %q on %s", op, bref)
		return nil
	}

	v, ok := types.NumConst(bref.T)
	if ok {
		switch op {
		case "+":
			return B // just shortcut this
		case "-":
			// a potential overflow here
			return tast.NewConst(tast.NewRef(types.NewNumber(-v)))
		}
		b.CodeErrorf(opTok.Pos, "pl.invalidOp",
			"invalid operation on num const: %q on %s", op, B)
		return nil
	}

	ct, ok := bref.T.(*types.Const)
	if !ok {
		b.CodeErrorf(opTok.Pos, "pl.expectConstExpr",
			"expect const expression but got %q %s", op, B)
		return nil
	}
	t := ct.Type
	v = ct.Value.(int64)
	if types.IsBasic(t, types.Bool) {
		// TODO
		b.CodeErrorf(opTok.Pos, "pl.notYetSupported",
			"const bool is not supported yet")
		return nil
	}
	if types.IsBasic(t, types.Float32) {
		b.CodeErrorf(opTok.Pos, "pl.notYetSupported",
			"float is not supported yet")
		return nil
	}
	if types.IsInteger(t) {
		switch op {
		case "+":
			return B // just shortcut this
		case "-":
			ref, e := types.NewConstInt(-v, t)
			if e != nil {
				b.CodeErrorf(opTok.Pos, "pl.constOverflow",
					"const %d overflows %q", v, t)
				return nil
			}
			return tast.NewConst(tast.NewRef(ref))
		}
		b.CodeErrorf(opTok.Pos, "pl.invalidOp",
			"invalid operation on int const: %q on %s", op, B)
		return nil
	}

	b.CodeErrorf(opTok.Pos, "pl.impossible",
		"only basic type const supported")
	return nil
}

func binaryOpConst(b *builder, opTok *lexing.Token, A, B tast.Expr) tast.Expr {
	op := opTok.Lit
	aref := A.R()
	bref := B.R()
	if !(aref.IsSingle() && bref.IsSingle()) {
		b.CodeErrorf(opTok.Pos, "pl.notSingle",
			"invalid expression list: %s %s %s", aref, op, bref)
		return nil
	}
	atyp := aref.Type()
	btyp := bref.Type()
	ca, oka := atyp.(*types.Const)
	cb, okb := btyp.(*types.Const)
	if !(oka && okb) {
		b.CodeErrorf(opTok.Pos, "pl.expectConstExpr",
			"expect a const expression on %s %s %s", atyp, op, btyp)
		return nil
	}

	va, oka := types.NumConst(atyp)
	vb, okb := types.NumConst(btyp)
	if oka && okb {
		return constIntOp(b, opTok, A, B, va, vb, ca.Type)
	}
	var t types.T
	if oka || okb {
		if oka && types.IsInteger(cb.Type) {
			if !types.InRange(va, cb.Type) {
				b.CodeErrorf(opTok.Pos, "pl.cannotCast",
					"cannot cast number %d to %s, out of range", va, cb.Type)
			}
			vb = cb.Value.(int64)
			t = cb.Type
		}
		if okb && types.IsInteger(ca.Type) {
			if !types.InRange(vb, ca.Type) {
				b.CodeErrorf(opTok.Pos, "pl.cannotCast",
					"cannot cast number %d to %s, out of range", vb, ca.Type)
			}
			va = ca.Value.(int64)
			t = ca.Type
		}
	} else {
		ta, oka := ca.Type.(types.Basic)
		tb, okb := cb.Type.(types.Basic)
		if !(oka && okb) {
			b.CodeErrorf(opTok.Pos, "pl.impossible",
				"only basic type const supported")
			return nil
		}

		if !(types.IsInteger(ta) && types.IsInteger(tb)) {
			b.CodeErrorf(opTok.Pos, "pl.notYetSupported",
				"only num and int consts are implemented")
			return nil
		}

		if tb != ta {
			b.CodeErrorf(
				opTok.Pos, "pl.invalidOp.typeMismatch",
				"cannot %s type %s, and type %s, type mismatch",
				op, ta, tb)
			return nil
		}
		va = ca.Value.(int64)
		vb = cb.Value.(int64)
		t = ta
	}
	return constIntOp(b, opTok, A, B, va, vb, t)
}

// TODO: after added const bool, remove inputs of va, ab
func constIntOp(b *builder, opTok *lexing.Token, A, B tast.Expr,
	va, vb int64, t types.T) tast.Expr {
	r := func(v int64) tast.Expr {
		if types.IsInteger(t) {
			ref, e := types.NewConstInt(v, t)
			if e != nil {
				b.CodeErrorf(opTok.Pos, "pl.constOverflow",
					"const %d overflows %q", v, t)
				return nil
			}
			return tast.NewConst(tast.NewRef(ref))
		}
		return tast.NewConst(tast.NewRef(types.NewNumber(v)))
	}
	op := opTok.Lit
	switch op {
	case "+":
		return r(va + vb)
	case "-":
		return r(va - vb)
	case "*":
		return r(va * vb)
	case "&":
		return r(va & vb)
	case "|":
		return r(va | vb)
	// remove ^ for nums?
	case "^":
		return r(va ^ vb)
	case "%":
		if vb == 0 {
			b.CodeErrorf(opTok.Pos, "pl.divideByZero", "modular by zero")
			return nil
		}
		return r(va % vb)
	case "/":
		if vb == 0 {
			b.CodeErrorf(opTok.Pos, "pl.divideByZero", "divide by zero")
			return nil
		}
		return r(va / vb)
	case "==", "!=", ">", "<", ">=", "<=":
		// TODO: will change into a const bool
		return &tast.OpExpr{
			A: A, Op: opTok, B: B,
			Ref: tast.NewRef(types.Bool),
		}
	case "<<":
		if vb < 0 {
			b.Errorf(opTok.Pos, "shift with negative value: %d", vb)
			return nil
		}
		return r(va << uint64(vb))
	case ">>":
		if vb < 0 {
			b.Errorf(opTok.Pos, "shift with negative value: %d", vb)
			return nil
		}
		return r(va >> uint64(vb))
	}

	b.CodeErrorf(opTok.Pos, "pl.invalidOp", "%q on int consts", op)
	return nil
}

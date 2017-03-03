package sempass

import (
	"math"
	"strconv"

	"shanhu.io/smlvm/lexing"
	"shanhu.io/smlvm/pl/ast"
	"shanhu.io/smlvm/pl/parse"
	"shanhu.io/smlvm/pl/tast"
	"shanhu.io/smlvm/pl/types"
)

func buildInt(b *builder, op *lexing.Token) tast.Expr {
	ret, e := strconv.ParseInt(op.Lit, 0, 64)
	if e != nil {
		b.CodeErrorf(op.Pos, "pl.cannotCast.invalidInteger",
			"invalid integer: %s", e)
		return nil
	}

	if ret > math.MaxUint32 {
		b.CodeErrorf(op.Pos, "pl.cannotCast.integerOverFlowed",
			"integer too large to fit in 32-bit")
		return nil
	}

	ref := tast.NewRef(types.NewNumber(ret))
	return tast.NewConst(ref)
}

func buildFloat(b *builder, op *lexing.Token) tast.Expr {
	ret, e := strconv.ParseFloat(op.Lit, 64)
	if e != nil {
		b.CodeErrorf(op.Pos, "pl.cannotCast.invalidFloat",
			"invalid float: %s", e)
		return nil
	}
	if ret != math.Floor(ret) {
		b.CodeErrorf(
			op.Pos, "pl.notYetSupported",
			"float is not yet supported")
		return nil
	}
	if ret > math.MaxUint32 {
		b.CodeErrorf(op.Pos, "pl.cannotCast.integerOverFlowed",
			"integer too large to fit in 32-bit")
		return nil
	}
	retInt := int64(ret)
	ref := tast.NewRef(types.NewNumber(retInt))
	return tast.NewConst(ref)
}

func buildChar(b *builder, op *lexing.Token) tast.Expr {
	v, e := strconv.Unquote(op.Lit)
	if e != nil {
		b.Errorf(op.Pos, "invalid char: %s", e)
		return nil
	} else if len(v) != 1 {
		b.Errorf(op.Pos, "invalid char in quote: %q", v)
		return nil
	}
	ref, e := types.NewConstInt(int64(v[0]), types.Int8)
	if e != nil {
		b.Errorf(op.Pos, "cannot covert char: %q in to int8", v)
		return nil
	}
	return tast.NewConst(tast.NewRef(ref))
}

func buildString(b *builder, op *lexing.Token) tast.Expr {
	v, e := strconv.Unquote(op.Lit)
	if e != nil {
		b.Errorf(op.Pos, "invalid string: %s", e)
		return nil
	}
	ref := tast.NewConstRef(types.String, v)
	return tast.NewConst(ref)
}

func buildIdent(b *builder, ident *lexing.Token) tast.Expr {
	s := b.scope.Query(ident.Lit)
	if s == nil {
		b.CodeErrorf(ident.Pos, "pl.undefinedIdent",
			"undefined identifier %s", ident.Lit)
		return nil
	}

	b.refSym(s, ident.Pos)

	t := s.ObjType.(types.T)
	switch s.Type {
	case tast.SymVar, tast.SymField:
		ref := tast.NewAddressableRef(t)
		return &tast.Ident{Token: ident, Ref: ref, Sym: s}
	case tast.SymConst:
		if types.IsConst(t) {
			return tast.NewConst(tast.NewRef(t))
		}
		ref := tast.NewRef(t)
		return &tast.Ident{Token: ident, Ref: ref, Sym: s}
	case tast.SymStruct, tast.SymType, tast.SymImport:
		ref := tast.NewRef(t)
		return &tast.Ident{Token: ident, Ref: ref, Sym: s}
	case tast.SymFunc:
		if t, ok := t.(*types.Func); ok {
			if t.MethodFunc == nil {
				return &tast.Ident{Token: ident, Ref: tast.NewRef(t), Sym: s}
			}
			if b.this == nil {
				panic("this missing")
			}
			ref := &tast.Ref{T: t.MethodFunc, Recv: b.this}
			return &tast.Ident{Token: ident, Ref: ref, Sym: s}
		}
		return &tast.Ident{Token: ident, Ref: tast.NewRef(t), Sym: s}
	default:
		b.Errorf(ident.Pos, "todo: token type: %s", tast.SymStr(s.Type))
		return nil
	}
}

func buildConstIdent(b *builder, ident *lexing.Token) tast.Expr {
	s := b.scope.Query(ident.Lit)
	if s == nil {
		b.CodeErrorf(ident.Pos, "pl.undefinedIdent",
			"undefined identifier %s", ident.Lit)
		return nil
	}

	b.refSym(s, ident.Pos)

	t := s.ObjType.(types.T)
	switch s.Type {
	case tast.SymConst:
		return tast.NewConst(tast.NewRef(t))
	case tast.SymStruct, tast.SymType, tast.SymImport:
		ref := tast.NewRef(t)
		return &tast.Ident{Token: ident, Ref: ref, Sym: s}
	}

	b.CodeErrorf(ident.Pos, "pl.expectConstExpr", "%s is a %s; expect a const",
		ident.Lit, tast.SymStr(s.Type),
	)
	return nil
}

func buildConstOperand(b *builder, op *ast.Operand) tast.Expr {
	switch op.Token.Type {
	case parse.Int:
		return buildInt(b, op.Token)
	case parse.Float:
		return buildFloat(b, op.Token)
	case parse.Ident:
		return buildConstIdent(b, op.Token)
	case parse.Char, parse.String:
		b.CodeErrorf(op.Token.Pos, "pl.notYetSupported",
			"only num and int consts are implemented")
		return nil
	}
	b.CodeErrorf(op.Token.Pos, "pl.expectConstExpr", "expect a constant")
	return nil
}

func buildOperand(b *builder, op *ast.Operand) tast.Expr {
	if op.Token.Type == parse.Keyword && op.Token.Lit == "this" {
		if b.this == nil {
			b.Errorf(op.Token.Pos, "using this out of a method function")
			return nil
		}
		return &tast.This{Ref: b.this}
	}

	switch op.Token.Type {
	case parse.Int:
		return buildInt(b, op.Token)
	case parse.Float:
		return buildFloat(b, op.Token)
	case parse.Char:
		return buildChar(b, op.Token)
	case parse.String:
		return buildString(b, op.Token)
	case parse.Ident:
		return buildIdent(b, op.Token)
	default:
		b.Errorf(op.Token.Pos, "invalid or not implemented: %d",
			op.Token.Type,
		)
		return nil
	}
}

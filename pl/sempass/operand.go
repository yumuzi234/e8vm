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
		b.Errorf(op.Pos, "invalid integer: %s", e)
		return nil
	}

	if ret < math.MinInt32 {
		b.Errorf(op.Pos, "integer too small to fit in 32-bit")
		return nil
	} else if ret > math.MaxUint32 {
		b.Errorf(op.Pos, "integer too large to fit in 32-bit")
		return nil
	}

	ref := tast.NewConstRef(types.NewNumber(ret), ret)
	return &tast.Const{ref}
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
	ref := tast.NewConstRef(types.Int8, int64(v[0]))
	return &tast.Const{ref}
}

func buildString(b *builder, op *lexing.Token) tast.Expr {
	v, e := strconv.Unquote(op.Lit)
	if e != nil {
		b.Errorf(op.Pos, "invalid string: %s", e)
		return nil
	}
	ref := tast.NewConstRef(types.String, v)
	return &tast.Const{ref}
}

func buildIdent(b *builder, ident *lexing.Token) tast.Expr {
	s := b.scope.Query(ident.Lit)
	if s == nil {
		b.Errorf(ident.Pos, "undefined identifier %s", ident.Lit)
		return nil
	}

	b.refSym(s, ident.Pos)

	t := s.ObjType.(types.T)
	switch s.Type {
	case tast.SymVar, tast.SymField:
		ref := tast.NewAddressableRef(t)
		return &tast.Ident{ident, ref, s}
	case tast.SymConst, tast.SymStruct, tast.SymType, tast.SymImport:
		ref := tast.NewRef(t)
		return &tast.Ident{ident, ref, s}
	case tast.SymFunc:
		if t, ok := t.(*types.Func); ok {
			if t.MethodFunc == nil {
				return &tast.Ident{ident, tast.NewRef(t), s}
			}
			if b.this == nil {
				panic("this missing")
			}
			ref := &tast.Ref{T: t.MethodFunc, Recv: b.this}
			return &tast.Ident{ident, ref, s}
		}
		return &tast.Ident{ident, tast.NewRef(t), s}
	default:
		b.Errorf(ident.Pos, "todo: token type: %s", tast.SymStr(s.Type))
		return nil
	}
}

func buildConstIdent(b *builder, ident *lexing.Token) tast.Expr {
	s := b.scope.Query(ident.Lit)
	if s == nil {
		b.Errorf(ident.Pos, "undefined identifier %s", ident.Lit)
		return nil
	}

	b.refSym(s, ident.Pos)

	t := s.ObjType.(types.T)
	switch s.Type {
	case tast.SymConst:
		ref := tast.NewRef(t)
		return &tast.Const{ref}
	case tast.SymStruct, tast.SymType, tast.SymImport:
		ref := tast.NewRef(t)
		return &tast.Ident{ident, ref, s}
	}

	b.Errorf(ident.Pos, "%s is a %s; expect a const",
		ident.Lit, tast.SymStr(s.Type),
	)
	return nil
}

func buildConstOperand(b *builder, op *ast.Operand) tast.Expr {
	switch op.Token.Type {
	case parse.Int:
		return buildInt(b, op.Token)
	case parse.Ident:
		return buildConstIdent(b, op.Token)
	}

	b.Errorf(op.Token.Pos, "expect a constant")
	return nil
}

func buildOperand(b *builder, op *ast.Operand) tast.Expr {
	if op.Token.Type == parse.Keyword && op.Token.Lit == "this" {
		if b.this == nil {
			b.Errorf(op.Token.Pos, "using this out of a method function")
			return nil
		}
		return &tast.This{b.this}
	}

	switch op.Token.Type {
	case parse.Int:
		return buildInt(b, op.Token)
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

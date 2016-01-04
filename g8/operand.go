package g8

import (
	"fmt"

	"e8vm.io/e8vm/g8/ir"
	"e8vm.io/e8vm/g8/tast"
	"e8vm.io/e8vm/g8/types"
)

func buildConst(b *builder, c *tast.Const) *ref {
	if _, ok := types.NumConst(c.T); ok {
		// untyped consts are just types.
		return newRef(c.T, nil)
	}

	if t, ok := c.T.(types.Basic); ok {
		v := c.ConstValue.(int64)
		switch t {
		case types.Int, types.Uint:
			return newRef(c.T, ir.Num(uint32(v)))
		case types.Int8, types.Uint8, types.Bool:
			return newRef(c.T, ir.Byt(uint8(v)))
		default:
			panic("other basic types not supported yet")
		}
	}

	if c.T == types.String {
		s := c.ConstValue.(string)
		ret := b.newTemp(c.T)
		b.b.Arith(ret.IR(), nil, "makeStr", b.p.NewString(s))
		return ret
	}

	panic("other const types not supported")
}

func buildField(b *builder, this ir.Ref, field *types.Field) *ref {
	retIR := ir.NewAddrRef(
		this,
		field.T.Size(),
		field.Offset(),
		types.IsByte(field.T),
		true,
	)
	return newAddressableRef(field.T, retIR)
}

func buildIdent(b *builder, id *tast.Ident) *ref {
	s := id.Symbol
	switch s.Type {
	case tast.SymVar:
		return s.Obj.(*objVar).ref
	case tast.SymFunc:
		v := s.Obj.(*objFunc)
		if !v.isMethod {
			return v.ref
		}
		if b.this == nil {
			panic("this missing")
		}
		return newRecvRef(v.Type().(*types.Func), b.this, v.IR())
	case tast.SymConst:
		return s.Obj.(*objConst).ref
	case tast.SymField:
		v := s.Obj.(*types.Field)
		return buildField(b, b.this.IR(), v)
	case tast.SymImport:
		return s.Obj.(*objImport).ref
	}
	panic(fmt.Errorf("unhandled token type: %s", tast.SymStr(s.Type)))
}

func buildConstIdent(b *builder, id *tast.Ident) *ref {
	s := id.Symbol
	switch s.Type {
	case tast.SymConst:
		return s.Obj.(*objConst).ref
	case tast.SymImport:
		return s.Obj.(*objImport).ref
	}
	panic(fmt.Errorf("not a const: %s", tast.SymStr(s.Type)))
}

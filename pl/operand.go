package pl

import (
	"fmt"

	"shanhu.io/smlvm/pl/codegen"
	"shanhu.io/smlvm/pl/tast"
	"shanhu.io/smlvm/pl/types"
)

func buildConst(b *builder, c *tast.Const) *ref {
	if _, ok := types.NumConst(c.T); ok {
		// untyped consts are just types.
		return newRef(c.T, nil)
	}

	if ct, ok := c.T.(*types.Const); ok {
		// typed consts
		return newRef(ct.Type, constNumIr(ct.Value.(int64), ct.Type))
	}

	if c.T == types.String {
		s := c.ConstValue.(string)
		ret := b.newTemp(c.T)
		b.b.Arith(ret.IR(), nil, "makeStr", b.p.NewString(s))
		return ret
	}

	// build slice literal
	if t, ok := c.T.(*types.Slice); ok {
		if bt, ok := t.T.(types.Basic); ok {
			switch bt {
			case types.Int, types.Uint, types.Int8, types.Uint8, types.Bool:
				bs := c.ConstValue.([]byte)
				ret := b.newTemp(t)
				ref := b.p.NewHeapDat(bs, bt.Size(), bt.RegSizeAlign())
				b.b.Arith(ret.IR(), nil, "makeDat", ref)
				return ret
			default:
				panic("other const slices not supported")
			}
		} else {
			panic("not basic type")
		}
	}

	panic("other const types not supported")
}

func buildField(b *builder, this codegen.Ref, field *types.Field) *ref {
	retIR := codegen.NewAddrRef(
		this,
		field.T.Size(),
		field.Offset(),
		types.IsByte(field.T),
		true,
	)
	return newAddressableRef(field.T, retIR)
}

func buildIdent(b *builder, id *tast.Ident) *ref {
	s := id.Sym
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
		t := s.ObjType.(types.T)
		return newRef(t, nil)
	}
	panic(fmt.Errorf("unhandled token type: %s", tast.SymStr(s.Type)))
}

func buildConstIdent(b *builder, id *tast.Ident) *ref {
	s := id.Sym
	switch s.Type {
	case tast.SymConst:
		return s.Obj.(*objConst).ref
	}
	panic(fmt.Errorf("not a const: %s", tast.SymStr(s.Type)))
}

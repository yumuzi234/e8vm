package sempass

import (
	"bytes"
	"encoding/binary"

	"e8vm.io/e8vm/g8/ast"
	"e8vm.io/e8vm/g8/tast"
	"e8vm.io/e8vm/g8/types"
)

func buildArrayLit(b *builder, lit *ast.ArrayLiteral) tast.Expr {
	buf := new(bytes.Buffer)

	t := buildType(b, lit.Type.Type)
	if t == nil {
		return nil
	}

	if !types.IsInteger(t) {
		pos := ast.ExprPos(lit.Type.Type)
		b.Errorf(pos, "array literal must be integer type")
		return nil
	}
	bt := t.(types.Basic)

	for _, expr := range lit.Exprs.Exprs {
		n := b.buildConstExpr(expr)
		if n == nil {
			return nil
		}
		ntype := n.R().T
		if _, ok := ntype.(*types.Const); !ok {
			b.Errorf(ast.ExprPos(expr), "array literal not a constant")
			return nil
		}

		if v, ok := types.NumConst(ntype); ok {
			if !types.InRange(v, t) {
				b.Errorf(ast.ExprPos(expr), "constant out of range of %s", t)
				return nil
			}

			switch bt {
			case types.Int, types.Uint:
				var bs [4]byte
				binary.LittleEndian.PutUint32(bs[:], uint32(v))
				buf.Write(bs[:])
			case types.Int8, types.Uint8:
				buf.Write([]byte{byte(v)})
			default:
				panic("not integer")
			}
		}
	}

	ref := tast.NewConstRef(&types.Slice{bt}, buf.Bytes())
	return &tast.Const{ref}
}

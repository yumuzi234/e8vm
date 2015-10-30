package g8

import (
	"e8vm.io/e8vm/g8/ast"
	"e8vm.io/e8vm/g8/types"
)

func buildFuncExit(b *builder) {
	next := b.f.NewBlock(b.b)
	b.b.Jump(b.f.End())
	b.b = next
}

func buildReturnStmt(b *builder, stmt *ast.ReturnStmt) {
	pos := stmt.Kw.Pos
	if stmt.Exprs == nil {
		if b.fretRef == nil || b.fretNamed {
			buildFuncExit(b)
		} else {
			b.Errorf(pos, "expects return %s", b.fretRef)
		}
	} else {
		if b.fretRef == nil {
			b.Errorf(pos, "function expects no return value")
			return
		}

		ref := buildExprList(b, stmt.Exprs)
		if ref == nil {
			return
		}

		nret := b.fretRef.Len()
		nsrc := ref.Len()
		if nret != nsrc {
			b.Errorf(pos, "expect (%s), returning (%s)", b.fretRef, ref)
			return
		}

		for i, t := range b.fretRef.typ {
			if !b.fretRef.addressable[i] {
				panic("bug")
			}

			srcType := ref.typ[i]
			if !types.CanAssign(t, srcType) {
				b.Errorf(pos, "expect (%s), returning (%s)", b.fretRef, ref)
				return
			}
		}

		if len(ref.addressable) > 1 {
			for i, addressable := range ref.addressable {
				if addressable {
					tmp := b.newTemp(ref.typ[i])
					b.b.Assign(tmp.IR(), ref.ir[i])
					ref.ir[i] = tmp.IR()
				}
			}
		}

		for i, dest := range b.fretRef.ir {
			if types.IsNil(ref.typ[i]) {
				b.b.Zero(dest)
			} else if v, ok := types.NumConst(ref.typ[i]); ok {
				b.b.Assign(dest, constNumIr(v, b.fretRef.typ[i]))
			} else {
				b.b.Assign(dest, ref.ir[i])
			}
		}

		buildFuncExit(b)
	}
}

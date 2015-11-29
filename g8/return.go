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

		src := buildExprList(b, stmt.Exprs)
		if src == nil {
			return
		}

		nret := b.fretRef.Len()
		nsrc := src.Len()
		if nret != nsrc {
			b.Errorf(pos, "expect (%s), returning (%s)", b.fretRef, src)
			return
		}

		for i := 0; i < nret; i++ {
			r := b.fretRef.At(i)
			if !r.Addressable() {
				panic("bug")
			}
			t := r.Type()

			srcType := src.At(i).Type()
			if !types.CanAssign(t, srcType) {
				b.Errorf(pos, "expect (%s), returning (%s)", b.fretRef, src)
				return
			}
		}

		if nsrc > 1 {
			srcCasted := new(ref)
			for i := 0; i < nsrc; i++ {
				r := src.At(i)
				if r.Addressable() {
					tmp := b.newTemp(r.Type())
					b.b.Assign(tmp.IR(), r.IR())
					srcCasted = appendRef(srcCasted, tmp)
				} else {
					srcCasted = appendRef(srcCasted, r)
				}
			}
			src = srcCasted
		}

		for i := 0; i < nret; i++ {
			r := b.fretRef.At(i)
			dest := r.IR()
			srcRef := src.At(i)
			t := srcRef.Type()
			if types.IsNil(t) {
				b.b.Zero(dest)
			} else if v, ok := types.NumConst(t); ok {
				b.b.Assign(dest, constNumIr(v, r.Type()))
			} else {
				b.b.Assign(dest, srcRef.IR())
			}
		}

		buildFuncExit(b)
	}
}

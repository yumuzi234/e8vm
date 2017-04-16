package sempass

import (
	"shanhu.io/smlvm/fmtutil"
	"shanhu.io/smlvm/pl/ast"
	"shanhu.io/smlvm/pl/tast"
)

func buildReturnStmt(b *builder, stmt *ast.ReturnStmt) tast.Stmt {
	pos := stmt.Kw.Pos
	if stmt.Exprs == nil {
		if b.retType == nil || b.retNamed {
			return &tast.ReturnStmt{}
		}
		b.CodeErrorf(pos, "pl.return.noReturnValue",
			"expects return %s", fmtutil.Join(b.retType, ","))
		return nil
	}

	if b.retType == nil {
		b.CodeErrorf(pos, "pl.return.expectNoReturn",
			"function expects no return value")
		return nil
	}

	src := b.buildExpr(stmt.Exprs)
	if src == nil {
		return nil
	}

	srcRef := src.R()
	nret := len(b.retType)
	nsrc := srcRef.Len()
	if nret != nsrc {
		b.CodeErrorf(pos, "pl.return.typeMismatch",
			"expect (%s), returning (%s)",
			fmtutil.Join(b.retType, ","), srcRef,
		)
		return nil
	}

	seenError := false
	cast := false
	mask := make([]bool, nret)
	expectRef := tast.Void

	for i := 0; i < nret; i++ {
		t := b.retType[i]
		srcType := srcRef.At(i).Type()
		ok, needCast := canAssign(b, pos, t, srcType)
		if !ok {
			b.CodeErrorf(pos, "pl.return.typeMismatch",
				"expect (%s), returning (%s)",
				fmtutil.Join(b.retType, ","), srcRef,
			)
		}
		seenError = seenError || !ok
		expectRef = tast.AppendRef(expectRef, tast.NewRef(t))
		cast = cast || needCast
		mask[i] = needCast
	}
	if seenError {
		return nil
	}

	// insert implicit type casts
	if cast {
		src = tast.NewMultiCast(src, expectRef, mask)
	}

	return &tast.ReturnStmt{Exprs: src}
}

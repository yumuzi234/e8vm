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

	srcTypes := srcRef.TypeList()
	res := canAssigns(b, pos, b.retType, srcTypes)
	if res.err {
		return nil
	}
	if res.needCast {
		src = tast.NewMultiCastTypes(src, b.retType, res.castMask)
	}

	return &tast.ReturnStmt{Exprs: src}
}

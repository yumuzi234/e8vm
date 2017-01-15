package sempass

import (
	"shanhu.io/smlvm/fmtutil"
	"shanhu.io/smlvm/pl/ast"
	"shanhu.io/smlvm/pl/tast"
	"shanhu.io/smlvm/pl/types"
)

func buildReturnStmt(b *builder, stmt *ast.ReturnStmt) tast.Stmt {
	pos := stmt.Kw.Pos
	if stmt.Exprs == nil {
		if b.retType == nil || b.retNamed {
			return &tast.ReturnStmt{}
		}
		b.Errorf(pos, "expects return %s", fmtutil.Join(b.retType, ","))
		return nil
	}

	if b.retType == nil {
		b.Errorf(pos, "function expects no return value")
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
		b.Errorf(pos, "expect (%s), returning (%s)",
			fmtutil.Join(b.retType, ","), srcRef,
		)
		return nil
	}

	for i := 0; i < nret; i++ {
		t := b.retType[i]
		srcType := srcRef.At(i).Type()
		if !types.CanAssign(t, srcType) {
			b.Errorf(pos, "expect (%s), returning (%s)",
				fmtutil.Join(b.retType, ","), srcRef,
			)
			return nil
		}
	}

	// insert implicit type casts
	if srcList, ok := tast.MakeExprList(src); ok {
		newList := tast.NewExprList()
		for i, e := range srcList.Exprs {
			t := e.Type()
			if types.IsNil(t) {
				e = tast.NewCast(e, b.retType[i])
			} else if v, ok := types.NumConst(t); ok {
				e = constCast(b, nil, v, e, b.retType[i])
				if e == nil {
					panic("bug")
				}
			}
			newList.Append(e)
		}
		src = newList
	}

	return &tast.ReturnStmt{Exprs: src}
}

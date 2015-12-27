package sempass

import (
	"e8vm.io/e8vm/g8/ast"
	"e8vm.io/e8vm/g8/tast"
	"e8vm.io/e8vm/g8/types"
	"e8vm.io/e8vm/lex8"
)

func assign(b *Builder, dest, src tast.Expr, op *lex8.Token) tast.Stmt {
	destRef := dest.R()
	srcRef := src.R()

	ndest := destRef.Len()
	nsrc := srcRef.Len()
	if ndest != nsrc {
		b.Errorf(op.Pos, "cannot assign %s to %s", nsrc, ndest)
		return nil
	}

	for i := 0; i < ndest; i++ {
		r := destRef.At(i)
		if !r.Addressable {
			b.Errorf(op.Pos, "assigning to non-addressable")
			return nil
		}

		destType := r.Type()
		srcType := srcRef.At(i).Type()
		if !types.CanAssign(destType, srcType) {
			b.Errorf(op.Pos, "cannot assign %s to %s", srcType, destType)
			return nil
		}
	}

	// insert casting if needed
	if srcList, ok := tast.MakeExprList(src); ok {
		newList := tast.NewExprList()
		for i, e := range srcList.Exprs {
			t := e.Type()
			if types.IsNil(t) {
				e = tast.NewCast(e, destRef.At(i).Type())
			} else if v, ok := types.NumConst(t); ok {
				e = constCast(b, nil, v, e, destRef.At(i).Type())
				if e == nil {
					panic("bug")
				}
			}
			newList.Append(e)
		}
		src = srcList
	}

	return &tast.AssignStmt{dest, src}
}

func opAssign(b *Builder, dest, src tast.Expr, op *lex8.Token) tast.Stmt {
	panic("todo")
}

func buildAssignStmt(b *Builder, stmt *ast.AssignStmt) tast.Stmt {
	left := b.BuildExpr(stmt.Left)
	if left == nil {
		return nil
	}

	right := b.BuildExpr(stmt.Right)
	if right == nil {
		return nil
	}

	if stmt.Assign.Lit == "=" {
		return assign(b, left, right, stmt.Assign)
	}

	return opAssign(b, left, right, stmt.Assign)
}

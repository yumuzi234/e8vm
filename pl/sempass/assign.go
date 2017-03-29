package sempass

import (
	"shanhu.io/smlvm/lexing"
	"shanhu.io/smlvm/pl/ast"
	"shanhu.io/smlvm/pl/tast"
	"shanhu.io/smlvm/pl/types"
)

func assign(b *builder, dest, src tast.Expr, op *lexing.Token) tast.Stmt {
	destRef := dest.R()
	srcRef := src.R()
	isError := false
	ndest := destRef.Len()
	nsrc := srcRef.Len()
	if ndest != nsrc {
		b.CodeErrorf(op.Pos, "pl.cannotAssign.lengthMismatch",
			"cannot assign(len) %s to %s; length mismatch",
			nsrc, ndest)
		return nil
	}

	for i := 0; i < ndest; i++ {
		r := destRef.At(i)
		if !r.Addressable {
			b.CodeErrorf(op.Pos, "pl.cannotAssign.notAddressable",
				"assigning to non-addressable")
			return nil
		}

		destType := r.Type()
		srcType := srcRef.At(i).Type()

		// assign for interface

		if !canAssign(b, op.Pos, destType, srcType) {
			isError = true
		}
	}
	if isError {
		return nil
	}

	// insert casting if needed
	if srcList, ok := tast.MakeExprList(src); ok {
		newList := tast.NewExprList()
		for i, e := range srcList.Exprs {
			t := e.Type()
			if types.IsNil(t) {
				e = tast.NewCast(e, destRef.At(i).Type())
			} else if v, ok := types.NumConst(t); ok {
				e = numCast(b, nil, v, e, destRef.At(i).Type())
				if e == nil {
					panic("bug")
				}
			}
			newList.Append(e)
		}
		src = newList
	}

	return &tast.AssignStmt{Left: dest, Op: op, Right: src}
}

func parseAssignOp(op string) string {
	opLen := len(op)
	if opLen == 0 {
		panic("invalid assign op")
	}
	return op[:opLen-1]
}

func opAssign(b *builder, dest, src tast.Expr, op *lexing.Token) tast.Stmt {
	destRef := dest.R()
	srcRef := src.R()
	if !destRef.IsSingle() || !srcRef.IsSingle() {
		b.CodeErrorf(op.Pos, "pl.cannotAssign.notSingle",
			"cannot assign %s %s %s", destRef, op.Lit, srcRef)
		return nil
	} else if !destRef.Addressable {
		b.CodeErrorf(op.Pos, "pl.cannotAssign.notAddressable",
			"assign to non-addressable")
		return nil
	}

	opLit := parseAssignOp(op.Lit)
	destType := destRef.Type()
	srcType := srcRef.Type()

	if opLit == ">>" || opLit == "<<" {
		if v, ok := types.NumConst(srcType); ok {
			src = numCast(b, op.Pos, v, src, types.Uint)
			if src == nil {
				return nil
			}
			srcRef = src.R()
			srcType = types.Uint
		}

		if !canShift(b, destType, srcType, op.Pos, opLit) {
			return nil
		}
		return &tast.AssignStmt{Left: dest, Op: op, Right: src}
	}

	if v, ok := types.NumConst(srcType); ok {
		src = numCast(b, op.Pos, v, src, destType)
		if src == nil {
			return nil
		}
		srcRef = src.R()
		srcType = destType
	}

	if ok, t := types.SameBasic(destType, srcType); ok {
		switch t {
		case types.Int, types.Int8, types.Uint, types.Uint8:
			return &tast.AssignStmt{Left: dest, Op: op, Right: src}
		}
	}

	b.Errorf(op.Pos, "invalid %s %s %s", destType, opLit, srcType)
	return nil
}

func buildAssignStmt(b *builder, stmt *ast.AssignStmt) tast.Stmt {
	hold := b.lhsSwap(true)
	left := b.buildExpr(stmt.Left)
	b.lhsRestore(hold)
	if left == nil {
		return nil
	}

	right := b.buildExpr(stmt.Right)
	if right == nil {
		return nil
	}

	if stmt.Assign.Lit == "=" {
		return assign(b, left, right, stmt.Assign)
	}

	return opAssign(b, left, right, stmt.Assign)
}

func canAssign(b *builder, p *lexing.Pos,
	left types.T, right types.T) bool {
	if i, ok := left.(*types.Interface); ok {
		// TODO
		if _, ok = right.(*types.Interface); ok {
			b.CodeErrorf(p, "pl.notYetSupported",
				"assign interface by interface is not supported yet")
			return false
		}
		if !assignInterface(b, p, i, right) {
			return false
		}
		return true
	}
	if !types.CanAssign(left, right) {
		b.CodeErrorf(p, "pl.cannotAssign.typeMismatch",
			"cannot assign %s to %s", left, right)
		return false
	}
	return true
}

func assignInterface(b *builder, p *lexing.Pos,
	i *types.Interface, right types.T) bool {
	flag := true
	s, ok := types.PointerOf(right).(*types.Struct)
	if !ok {
		b.CodeErrorf(p, "pl.cannotAssign.interface",
			"cannot assign interface %s by %s, not a struct pointer", i, right)
		return false
	}
	e := func(fname, m string) {
		b.CodeErrorf(p, "pl.cannotAssign.interface",
			"cannot assign interface %s by %s, "+m, i, right, fname)
	}

	funcs := i.Syms.List()
	for _, f := range funcs {
		sym := s.Syms.Query(f.Name())
		if sym == nil {
			e(f.Name(), "func %s not implemented")
			flag = false
			continue
		}
		t2, ok := sym.ObjType.(*types.Func)
		if !ok {
			e(f.Name(), "%s is a struct member but not a method")
			flag = false
			continue
		}
		t2 = t2.MethodFunc
		t1 := f.ObjType.(*types.Func)
		if len(t1.Args) != len(t2.Args) {
			e(f.Name(), "args number mismatch for %s")
			flag = false
			continue
		}
		if len(t1.Rets) != len(t2.Rets) {
			e(f.Name(), "returns number mismatch for %s")
			flag = false
			continue
		}
		for i, t := range t1.Args {
			if !types.SameType(t.T, t2.Args[i].T) {
				b.CodeErrorf(p, "pl.cannotAssign.interface",
					"cannot assign interface %s by %s, "+
						"type not match, %v, %v in func %s",
					i, right, t.T, t2.Args[i].T, f.Name())
				flag = false
				continue
			}
		}

		for i, t := range t1.Rets {
			if !types.SameType(t.T, t2.Rets[i].T) {
				b.CodeErrorf(p, "pl.cannotAssign.interface",
					"cannot assign interface %s by %s, "+
						"type not match, %v, %v in func %s",
					i, right, t.T, t2.Args[i].T, f.Name())
				flag = false
				continue
			}
		}
	}
	return flag
}

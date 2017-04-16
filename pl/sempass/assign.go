package sempass

import (
	"fmt"

	"shanhu.io/smlvm/lexing"
	"shanhu.io/smlvm/pl/ast"
	"shanhu.io/smlvm/pl/tast"
	"shanhu.io/smlvm/pl/types"
)

func assign(b *builder, dest, src tast.Expr, op *lexing.Token) tast.Stmt {
	destRef := dest.R()
	srcRef := src.R()
	ndest := destRef.Len()
	nsrc := srcRef.Len()
	if ndest != nsrc {
		b.CodeErrorf(op.Pos, "pl.cannotAssign.lengthMismatch",
			"cannot assign(len) %s to %s; length mismatch",
			nsrc, ndest)
		return nil
	}

	// check if all addressable
	for i := 0; i < ndest; i++ {
		r := destRef.At(i)
		if !r.Addressable {
			b.CodeErrorf(
				op.Pos, "pl.cannotAssign.notAddressable",
				"assigning to non-addressable",
			)
			return nil
		}
	}

	seenError := false
	cast := false
	mask := make([]bool, ndest)

	for i := 0; i < ndest; i++ {
		r := destRef.At(i)
		destType := r.Type()
		srcType := srcRef.At(i).Type()

		// assign for interface
		ok, needCast := canAssign(b, op.Pos, destType, srcType)
		seenError = seenError || !ok
		cast = cast || needCast
		mask[i] = needCast
	}
	if seenError {
		return nil
	}

	// insert casting if needed
	if cast {
		src = tast.NewMultiCast(src, destRef, mask)
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

func canAssign(
	b *builder, p *lexing.Pos, left, right types.T,
) (ok bool, needCast bool) {
	if i, ok := left.(*types.Interface); ok {
		// TODO(yumuzi234): assing interface from interface
		if _, ok = right.(*types.Interface); ok {
			b.CodeErrorf(p, "pl.notYetSupported",
				"assign interface by interface is not supported yet")
			return false, false
		}
		if !assignInterface(b, p, i, right) {
			return false, false
		}
		return true, true
	}
	ok, needCast = types.CanAssign(left, right)
	if !ok {
		b.CodeErrorf(p, "pl.cannotAssign.typeMismatch",
			"cannot use %s as %s", left, right)
		return false, false
	}
	return ok, needCast
}

func assignInterface(b *builder, p *lexing.Pos,
	i *types.Interface, right types.T) bool {
	flag := true
	s, ok := types.PointerOf(right).(*types.Struct)
	if !ok {
		b.CodeErrorf(p, "pl.cannotAssign.interface",
			"cannot use %s as interface %s, not a struct pointer", right, i)
		return false
	}
	errorf := func(f string, a ...interface{}) {
		m := fmt.Sprintf(f, a...)
		b.CodeErrorf(p, "pl.cannotAssign.interface",
			"cannot use %s as interface %s, %s", right, i, m)
		flag = false
	}

	funcs := i.Syms.List()
	for _, f := range funcs {
		sym := s.Syms.Query(f.Name())
		if sym == nil {
			errorf("function %s not implemented", f.Name())
			continue
		}
		t2, ok := sym.ObjType.(*types.Func)
		if !ok {
			errorf("%s is a struct member but not a method", f.Name())
			continue
		}
		t2 = t2.MethodFunc
		t1 := f.ObjType.(*types.Func)
		if !types.SameType(t1, t2) {
			errorf("func signature mismatch %q, %q", t1, t2)
		}
	}
	return flag
}

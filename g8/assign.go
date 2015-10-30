package g8

import (
	"e8vm.io/e8vm/g8/ast"
	"e8vm.io/e8vm/g8/types"
	"e8vm.io/e8vm/lex8"
)

func assign(b *builder, dest, src *ref, op *lex8.Token) bool {
	ndest := dest.Len()
	nsrc := src.Len()
	if ndest != nsrc {
		b.Errorf(op.Pos, "cannot assign %s to %s",
			nsrc, ndest,
		)
		return false
	}

	for i, destType := range dest.typ {
		if !dest.addressable[i] {
			b.Errorf(op.Pos, "assigning to non-addressable")
			return false
		}

		srcType := src.typ[i]
		if !types.CanAssign(destType, srcType) {
			b.Errorf(op.Pos, "cannot assign %s to %s", src, dest)
			return false
		}
	}

	if len(dest.ir) == 1 {
		srcTyp := src.typ[0]
		destTyp := dest.typ[0]
		destIr := dest.ir[0]
		if types.IsNil(srcTyp) {
			b.b.Zero(destIr)
		} else if v, ok := types.NumConst(srcTyp); ok {
			b.b.Assign(destIr, constNumIr(v, destTyp))
		} else {
			b.b.Assign(destIr, src.ir[0])
		}
	} else {
		temps := make([]*ref, len(dest.ir))
		// perform the assignment
		for i, srcIr := range src.ir {
			srcTyp := src.typ[i]
			if types.IsNil(srcTyp) {
				continue // will zero directly later
			}
			if _, ok := srcTyp.(*types.Const); ok {
				continue
			}

			temps[i] = b.newTemp(srcTyp)
			b.b.Assign(temps[i].IR(), srcIr)
		}

		for i, destIr := range dest.ir {
			srcTyp := src.typ[i]
			if types.IsNil(srcTyp) {
				b.b.Zero(destIr)
			} else if v, ok := types.NumConst(srcTyp); ok {
				b.b.Assign(destIr, constNumIr(v, dest.typ[i]))
			} else {
				b.b.Assign(destIr, temps[i].IR())
			}
		}
	}

	return true
}

func parseAssignOp(op string) string {
	opLen := len(op)
	if opLen == 0 {
		panic("invalid assign op op")
	}
	return op[:opLen-1]
}

func opAssignInt(b *builder, opOp string, dest, src *ref) {
	switch opOp {
	case "+", "-", "*", "&", "|", "^", "/", "%":
		buildBasicArith(b, dest, dest, src, opOp)
	}
}

func opAssign(b *builder, dest, src *ref, op *lex8.Token) {
	if !dest.IsSingle() || !src.IsSingle() {
		b.Errorf(op.Pos, "%s %s %s", dest, op.Lit, src)
		return
	} else if !dest.Addressable() {
		b.Errorf(op.Pos, "assign to non-addressable")
		return
	}

	opOp := parseAssignOp(op.Lit)

	destType := dest.Type()
	srcType := src.Type()

	if opOp == ">>" || opOp == "<<" {
		if v, ok := types.NumConst(srcType); ok {
			src = constCast(b, op.Pos, v, types.Uint)
			if src == nil {
				return
			}
			srcType = types.Uint
		}

		if canShift(b, destType, srcType, op.Pos, opOp) {
			buildShift(b, dest, dest, src, opOp)
		}
		return
	}

	if v, ok := types.NumConst(srcType); ok {
		src = constCast(b, op.Pos, v, destType)
		if src == nil {
			return
		}
		srcType = destType
	}

	if ok, t := types.SameBasic(destType, srcType); ok {
		switch t {
		case types.Int, types.Int8, types.Uint, types.Uint8:
			opAssignInt(b, opOp, dest, src)
			return
		}
	}

	b.Errorf(op.Pos, "invalid %q", op.Lit)
}

func buildAssignStmt(b *builder, stmt *ast.AssignStmt) {
	left := buildExprList(b, stmt.Left)
	if left == nil {
		return
	}
	right := buildExprList(b, stmt.Right)
	if right == nil {
		return
	}

	if stmt.Assign.Lit == "=" {
		assign(b, left, right, stmt.Assign)
		return
	}

	opAssign(b, left, right, stmt.Assign)
}

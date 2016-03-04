package g8

import (
	"e8vm.io/e8vm/g8/ir"
	"e8vm.io/e8vm/g8/tast"
	"e8vm.io/e8vm/g8/types"
)

func buildAssignStmt(b *builder, stmt *tast.AssignStmt) {
	left := b.buildExpr(stmt.Left)
	right := b.buildExpr(stmt.Right)
	if stmt.Op.Lit == "=" {
		assign(b, left, right)
		return
	}
	opAssign(b, left, right, stmt.Op.Lit)
}

func assign(b *builder, dest, src *ref) {
	n := dest.Len()
	if n == 1 {
		b.b.Assign(dest.IR(), src.IR())
		return
	}

	temps := make([]ir.Ref, n)
	for i := 0; i < n; i++ {
		s := src.At(i)
		t := s.Type()
		temps[i] = b.newTempIR(t)
		b.b.Assign(temps[i], s.IR())
	}

	for i := 0; i < n; i++ {
		b.b.Assign(dest.At(i).IR(), temps[i])
	}
}

func opAssign(b *builder, dest, src *ref, op string) {
	opOp := op[:len(op)-1]
	if opOp == ">>" || opOp == "<<" {
		buildShift(b, dest, dest, src, opOp)
		return
	}

	ok, t := types.SameBasic(src.Type(), dest.Type())
	if !ok {
		panic("bug")
	}

	switch t {
	case types.Int, types.Int8:
		buildBasicArith(b, dest, dest, src, opOp)
		return
	case types.Uint, types.Uint8:
		switch opOp {
		case "*", "/", "%":
			opOp = "u" + opOp
		}
		buildBasicArith(b, dest, dest, src, opOp)
		return
	}

	panic("bug")
}

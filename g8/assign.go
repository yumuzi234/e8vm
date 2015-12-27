package g8

import (
	"e8vm.io/e8vm/g8/ir"
	"e8vm.io/e8vm/g8/tast"
)

func buildAssignStmt(b *builder, stmt *tast.AssignStmt) {
	left := b.buildExpr2(stmt.Left)
	right := b.buildExpr2(stmt.Right)
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
	if opOp == ">>=" || opOp == "<<=" {
		buildShift(b, dest, dest, src, opOp)
		return
	}

	buildBasicArith(b, dest, dest, src, opOp)
}

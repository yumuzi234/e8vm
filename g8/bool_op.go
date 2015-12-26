package g8

import (
	"e8vm.io/e8vm/g8/tast"
	"e8vm.io/e8vm/g8/types"
)

func binaryOpBool(b *builder, op string, A, B *ref) *ref {
	switch op {
	case "==", "!=":
		ret := b.newTemp(types.Bool)
		b.b.Arith(ret.IR(), A.IR(), op, B.IR())
		return ret
	}
	panic("bug")
}

func buildLogicAnd(b *builder, A *ref, B tast.Expr) *ref {
	blockB := b.f.NewBlock(b.b)
	retFalse := b.f.NewBlock(blockB)
	after := b.f.NewBlock(retFalse)

	ret := b.newTemp(types.Bool)

	b.b.JumpIfNot(A.IR(), retFalse)

	// evaluate expression B
	b.b = blockB
	refB := b.buildExpr2(B)
	b.b.Assign(ret.IR(), refB.IR()) // and save it as result

	b.b.Jump(after)

	retFalse.Assign(ret.IR(), refFalse.IR())
	retFalse.Jump(after)

	b.b = after
	return ret
}

func buildLogicOr(b *builder, A *ref, B tast.Expr) *ref {
	blockB := b.f.NewBlock(b.b)
	retTrue := b.f.NewBlock(blockB)
	after := b.f.NewBlock(retTrue)

	ret := b.newTemp(types.Bool)

	b.b.JumpIf(A.IR(), retTrue)

	// evaluate expression B
	b.b = blockB
	refB := b.buildExpr2(B)
	b.b.Assign(ret.IR(), refB.IR()) // and save it as result
	b.b.Jump(after)

	retTrue.Assign(ret.IR(), refTrue.IR())
	retTrue.Jump(after)

	b.b = after
	return ret
}

func unaryOpBool(b *builder, op string, B *ref) *ref {
	if op != "!" {
		panic("bug")
	}

	ret := b.newTemp(types.Bool)
	b.b.Arith(ret.IR(), nil, op, B.IR())
	return ret
}

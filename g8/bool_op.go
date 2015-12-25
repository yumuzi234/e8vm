package g8

import (
	"e8vm.io/e8vm/g8/ast"
	"e8vm.io/e8vm/g8/types"
	"e8vm.io/e8vm/lex8"
)

func binaryOpBool(b *builder, opTok *lex8.Token, A, B *ref) *ref {
	op := opTok.Lit
	switch op {
	case "==", "!=":
		ret := b.newTemp(types.Bool)
		b.b.Arith(ret.IR(), A.IR(), op, B.IR())
		return ret
	}

	b.Errorf(opTok.Pos, "%q on bools", op)
	return nil
}

func buildBoolExprForOp(b *builder, opTok *lex8.Token, B ast.Expr) *ref {
	ret := b.buildExpr(B)
	if ret == nil {
		// error
	} else if !ret.IsSingle() {
		b.Errorf(opTok.Pos, "%q with %s", opTok.Lit, ret)
		ret = nil
	} else if !types.IsBasic(ret.Type(), types.Bool) {
		b.Errorf(opTok.Pos, "%q with %s", opTok.Lit, ret)
		ret = nil
	}
	return ret
}

func buildLogicAnd(b *builder, opTok *lex8.Token, A *ref, B ast.Expr) *ref {
	blockB := b.f.NewBlock(b.b)
	retFalse := b.f.NewBlock(blockB)
	after := b.f.NewBlock(retFalse)

	ret := b.newTemp(types.Bool)

	b.b.JumpIfNot(A.IR(), retFalse)

	// evaluate expression B
	b.b = blockB
	refB := buildBoolExprForOp(b, opTok, B)
	if refB != nil {
		b.b.Assign(ret.IR(), refB.IR()) // and save it as result
	}

	b.b.Jump(after)

	retFalse.Assign(ret.IR(), refFalse.IR())
	retFalse.Jump(after)

	b.b = after

	if refB == nil {
		return nil
	}
	return ret
}

func buildLogicOr(b *builder, opTok *lex8.Token, A *ref, B ast.Expr) *ref {
	blockB := b.f.NewBlock(b.b)
	retTrue := b.f.NewBlock(blockB)
	after := b.f.NewBlock(retTrue)

	ret := b.newTemp(types.Bool)

	b.b.JumpIf(A.IR(), retTrue)

	// evaluate expression B
	b.b = blockB
	refB := buildBoolExprForOp(b, opTok, B)
	if refB != nil {
		b.b.Assign(ret.IR(), refB.IR()) // and save it as result
	}
	b.b.Jump(after)

	retTrue.Assign(ret.IR(), refTrue.IR())
	retTrue.Jump(after)

	b.b = after

	if refB == nil {
		return nil
	}
	return ret
}

func unaryOpBool(b *builder, opTok *lex8.Token, B *ref) *ref {
	op := opTok.Lit
	switch op {
	case "!":
		ret := b.newTemp(types.Bool)
		b.b.Arith(ret.IR(), nil, op, B.IR())
		return ret
	}
	b.Errorf(opTok.Pos, "invalid operation: %q on boolean", op)
	return nil
}

package pl

import (
	"shanhu.io/smlvm/pl/codegen"
)

func nilPointerPanic(b *builder, pt codegen.Ref) {
	nilPointerPanicOp(b, pt, "?")
}

func nilFuncPointerPanic(b *builder, pt codegen.Ref) {
	nilPointerPanicOp(b, pt, "?f")
}

func nilPointerPanicOp(b *builder, pt codegen.Ref, op string) {
	// TODO: optimize for consts/functions that are impossible to be nil
	if !codegen.CanBeZero(pt) {
		return
	}

	nonZero := b.newCond()
	b.b.Arith(nonZero, nil, op, pt)

	body := b.f.NewBlock(b.b)
	after := b.f.NewBlock(body)

	b.b.JumpIf(nonZero, after)

	b.b = body
	callPanic(b, "reference nil pointer")

	b.b = after
}

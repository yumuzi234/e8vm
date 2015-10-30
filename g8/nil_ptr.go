package g8

import (
	"e8vm.io/e8vm/g8/ir"
)

func nilPointerPanic(b *builder, pt ir.Ref) {
	nilPointerPanicOp(b, pt, "?")
}

func nilFuncPointerPanic(b *builder, pt ir.Ref) {
	nilPointerPanicOp(b, pt, "?f")
}

func nilPointerPanicOp(b *builder, pt ir.Ref, op string) {
	// TODO: optimize for consts/functions that are impossible to be nil
	if !ir.CanBeZero(pt) {
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

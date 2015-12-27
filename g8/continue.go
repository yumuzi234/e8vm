package g8

func genContinueStmt(b *builder) {
	after := b.f.NewBlock(b.b)
	b.b.Jump(b.continues.top())
	b.b = after
}

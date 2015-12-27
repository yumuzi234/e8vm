package g8

func buildBreakStmt(b *builder) {
	after := b.f.NewBlock(b.b)
	b.b.Jump(b.breaks.top())
	b.b = after
}

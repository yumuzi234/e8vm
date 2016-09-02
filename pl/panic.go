package pl

func callPanic(b *builder, msg string) {
	if b.panicFunc == nil {
		panic("panic function missing")
	}
	// TODO: print a message
	b.b.Call(nil, b.panicFunc)
}

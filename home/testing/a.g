func testMain(f func()) {
	printInt(3333)
	f()
}

func TestBadSomething() { panic() }
func TestSomething() { }

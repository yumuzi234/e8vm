struct A { 
	a int;
	func s(a int) { this.a = a }
	func p() { printInt(a) }
}
func main() { var a A; a.s(33); a.p() }
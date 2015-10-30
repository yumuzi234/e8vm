struct A { a, b int; func p() { printInt(a+b) } }; var a A
func main() { a.a=30; a.b=3; a.p() }
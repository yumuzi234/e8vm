struct A { a int; func p() { printInt(a) } }; var a A
func main() { a.a=33; a.p() }
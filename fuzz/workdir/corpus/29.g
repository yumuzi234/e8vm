struct A { func p() { printInt(33) }; func q() { p() } }
func main() { var a A; a.q() }
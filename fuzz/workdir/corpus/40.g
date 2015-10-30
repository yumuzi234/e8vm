struct a { a int; b,c byte }
func main() { 
	var as [4]a
	printUint(uint(&as[1])-uint(&as[0]))
	printUint(uint(&as[0].c)-uint(&as[0].a))
}
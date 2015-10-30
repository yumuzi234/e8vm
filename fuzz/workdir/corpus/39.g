struct a { a int; b byte }
func main() { 
	var x,y a; x.a=33; y.a=44; 
	printInt(x.a); printInt(y.a) 
}
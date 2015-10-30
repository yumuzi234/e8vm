var ( a, b []int; v [3]int )
func main() { a=v[:]; b = v[:]; if a == b { printInt(33) } }
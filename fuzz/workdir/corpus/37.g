var ( a, b []int; v [3]int )
func main() { a=v[:2]; b = v[:3]; if a != b { printInt(33) } }
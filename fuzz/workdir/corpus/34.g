var ( a []int; s [2]int )
func main() { a=s[:]; s[1]=33; printInt(a[1]) }
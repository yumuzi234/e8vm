func f(a []int) { printInt(a[3]) }
func main() { var a [8]int; a[4]=33; f(a[1:5]) }
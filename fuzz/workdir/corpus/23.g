func p(i int) { printInt(i+2) }
func c(x func(i int)) { x(33) }
func main() { c(p) }
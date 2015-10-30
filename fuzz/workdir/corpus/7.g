func b() bool { printInt(4); return true }
func main() { if false || b() { printInt(3) } }

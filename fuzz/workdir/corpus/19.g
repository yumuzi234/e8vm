struct A { a int }
func main() { var a A; pi:=&a.a; *pi=33; printInt(a.a) }
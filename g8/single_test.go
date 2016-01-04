package g8

import (
	"testing"

	"errors"
	"strings"

	"e8vm.io/e8vm/arch8"
)

var errRunFailed = errors.New("test run failed")

func singleTestRun(t *testing.T, input string, N int) (string, error) {
	bs, es, _ := CompileSingle("main.g", input, false)
	if es != nil {
		t.Log(input)
		for _, err := range es {
			t.Log(err)
		}
		t.Error("compile failed")
		return "", errRunFailed
	}

	ncycle, out, err := arch8.RunImageOutput(bs, N)
	if ncycle == N {
		t.Log(input)
		t.Error("running out of time")
		return "", errRunFailed
	}
	return out, err
}

func TestSingleFile(t *testing.T) {
	const N = 100000

	o := func(input, output string) {
		out, err := singleTestRun(t, input, N)
		if err == errRunFailed {
			t.Error(err)
			return
		}
		if !arch8.IsHalt(err) {
			t.Log(input)
			t.Log(err)
			t.Error("did not halt gracefully")
			return
		}

		got := strings.TrimSpace(out)
		expect := strings.TrimSpace(output)
		if got != expect {
			t.Log(input)
			t.Logf("expect: %s", expect)
			t.Errorf("got: %s", got)
		}
	}

	o("func main() { }", "")
	o("func main() { return }", "")
	o("func main() { printInt(3) }", "3")

	o(` func r() int { return 7 }
		func main() { printInt(r()) }`,
		"7")
	o(`	func p(i int) { printInt(i) }
		func main() { p(33); p(44) }`,
		"33\n44")
	o(`	func r() (int, int) { return 3, 4 }
		func main() { a, b := r(); printInt(a); printInt(b) }`,
		"3\n4")
	o(` func r() (int, int) { return 3, 4 }
		func p(a, b int) { printInt(a); printInt(b) }
		func main() { p(r()) }`,
		"3\n4")
	o(` func r() (int) { return 3 }
		func p(a, b int) { printInt(a); printInt(b) }
		func main() { p(r(), 4) }`,
		"3\n4")
	o(`	func fabo(n int) int {
			if n == 0 return 0
			if n == 1 return 1
			return fabo(n-1) + fabo(n-2)
		}
		func main() { printInt(fabo(10)) }`,
		"55")

	o(`	func b() bool { printInt(4); return true }
		func main() { if false || b() { printInt(3) } }`,
		"4\n3")
	o(`	func b() bool { printInt(4); return true }
		func main() { if true || b() { printInt(3) } }`,
		"3")
	o(`	func b() bool { printInt(4); return false }
		func main() { if true && b() { printInt(3) } }`,
		"4")
	o(`	func b() bool { printInt(4); return true }
		func main() { if false && b() { printInt(3) } }`,
		"")
	o(`	func b() bool { printInt(4); return true }
		func main() { if true && b() { printInt(3) } }`,
		"4\n3")
	o(`func f(i int) { i=33; printInt(i) }; func main() { f(44) }`, "33")
	o(`	func f(a []int) { printInt(a[3]) }
		func main() { var a [8]int; a[4]=33; f(a[1:5]) }`,
		"33")
	o(`func main() { for true { printInt(33); break } }`, "33")
	o(`func main() { for 0==0 { printInt(33); break } }`, "33")
	o(`func main() { for true && true { printInt(33); break } }`, "33")
	o(`func main() { for true && false { printInt(33); break } }`, "")
	o(`func main() { for 0==0 && 0==1 { printInt(33); break } }`, "")
	o(`func main() { for true || false { printInt(33); break } }`, "33")
	o(`func main() { for false || true { printInt(33); break } }`, "33")
	o(`func main() { for false { printInt(33); break } }`, "")

	o(`	func printStr(s string) {
			n:=len(s); for i:=0;i<n;i++ { printChar(s[i]) }
		}
		func main() { printStr("hello") }`, "hello")
	o(`	func printStr(s string) {
			n:=len(s); for i:=0;i<n;i++ { printChar(s[i]) }
		}
		func main() { var a []int8; b:="hello"; a=b; printStr(a) }`, "hello")
	o(`	func printStr(s string) {
			n:=len(s); for i:=0;i<n;i++ { printChar(s[i]) }
		}
		func main() { b:="hello"; var a []int8; a=b; printStr(b) }`, "hello")

	o(` struct A {}; func main() { var a A; pa := &a; }`, "")
	o(` struct A { a int }
		func main() { var a A; printInt(a.a) }`, "0")
	o(` struct A { a int }
		func main() { var a A; a.a = 33; printInt(a.a) }`, "33")
	o(` struct A { a int }
		func main() { var a A; pi:=&a.a; *pi=33; printInt(a.a) }`, "33")
	o(` struct A { a int }
		func main() { var a A; (&a).a = 33; printInt(a.a) }`, "33")
	o(` struct A { a int }
		func main() { var a A; var pa=&a; pa.a = 33; printInt(pa.a) }`, "33")
	o(` struct A { b byte; a int }
		func main() { var a A; var pa=&a; 
			pa.a = 33; pa.b = byte(7);
			printInt(pa.a); printInt(int(pa.b))
		}`, "33\n7")

	o(` func p(i int) { printInt(i) }
		func main() { f:=p; f(33) }`, "33")
	o(` func p(i int) { printInt(i+2) }
		func c(x func(i int)) { x(33) }
		func main() { c(p) }`, "35")

	o("struct A { a *A }; func main() {}", "")
	o(`struct A { b B }; struct B { a *A }; func main() {}`, "")

	o(` struct A { func p(a int) { printInt(a) } }
		func main() { var a A; a.p(33) }`, "33")
	o(` struct A { 
			a int;
			func s(a int) { this.a = a }
			func p() { printInt(a) }
		}
		func main() { var a A; a.s(33); a.p() }`, "33")
	o(` struct A { 
			a int;
			func s(a int) { (*this).a = a }
			func p() { printInt(a) }
		}
		func main() { var a A; a.s(33); a.p() }`, "33")
	o(` struct A { 
			a int;
			func p() { printInt(a) }
			func q(a int) { printInt(a) }
		}
		func main() { var a A; a.p(); a.a=33; a.p(); a.q(78) }`, "0\n33\n78")
	o(` struct A { func p() { printInt(33) }; func q() { p() } }
		func main() { var a A; a.q() }`, "33")

	o("var a int; func main() { a := 33; printInt(a) }", "33")
	o(`	struct A { func p() { printInt(33) } }; var a A
		func main() { a.p() }`, "33")
	o(` struct A { a int; func p() { printInt(a) } }; var a A
		func main() { a.a=33; a.p() }`, "33")
	o(` struct A { a, b int; func p() { printInt(a+b) } }; var a A
		func main() { a.a=30; a.b=3; a.p() }`, "33")
	o(` var ( a []int; s [2]int )
		func main() { a=s[:]; s[1]=33; printInt(a[1]) }`, "33")
	o(` var ( a []int )
		func main() { if a == nil { printInt(33) } }`, "33")
	o(` var ( a []int; v [0]int )
		func main() { a=v[:]; if a != nil { printInt(33) } }`, "33")
	o(` var ( a, b []int; v [3]int )
		func main() { a=v[:2]; b = v[:3]; if a != b { printInt(33) } }`, "33")
	o(` var ( a, b []int; v [3]int )
		func main() { a=v[:]; b = v[:]; if a == b { printInt(33) } }`, "33")

	o("func init() { printInt(33) }; func main() { printInt(44) }", "33\n44")

	o(`	struct a { a int; b byte }
		func main() { 
			var x,y a; x.a=33; y.a=44; 
			printInt(x.a); printInt(y.a) 
		}`, "33\n44")
	o(` struct a { a int; b,c byte }
		func main() { 
			var as [4]a
			printUint(uint(&as[1])-uint(&as[0]))
			printUint(uint(&as[0].c)-uint(&as[0].a))
		}`, "8\n5")
	o(` struct a { a [4]byte; b byte }
		func main() { 
			var as [4]a
			printUint(uint(&as[1])-uint(&as[0]))
		}`, "5")

	o("import(); func main() {}", "")
	o("import(); func main() { printInt(33) }", "33")

	o("const a=33; func main() { printInt(a) }", "33")
	o("const a=33; func main() { const a=34; printInt(a) }", "34")
	o(`	const a, b = 1, a; const c, d = d, 3
		func main() { printInt(a); printInt(b); printInt(c); printInt(d) }`,
		"1\n1\n3\n3\n")
	o(`const a, b = b + 3, 30; func main() { printInt(a) }`, "33")

	o(`	var a [4]int
		func main() {
			s := make([]int, 2, &a[1])
			s[0] = 33; s[1] = 47
			printInt(a[1]); printInt(a[2])
		}`, "33\n47")

	// Bugs found by the fuzzer in the past
	o("func main() { a := 0==0; if a { printInt(33) } }", "33")
	o("func n()[(3+4)*5]string{}; func main() { printInt(len(n())) }", "35")
	o("func n()[1<<6]string{}; func main() { printInt(len(n())) }", "64")
	o("func n()[1>>6]string{}; func main() { printInt(len(n())) }", "0")
	o("func main() { r:=+'0'; printChar(r) }", "0")
	o("func n(b**bool) { **b=**b }; func main() {}", "")
	o("func n(b****bool) { ****b=****b }; func main() {}", "")
	o(` struct A { func n() (a int) { return 33 } };
		func main() { var a A; printInt(a.n()) }`, "33")
	o(`	func main() {printInt(33)}
		func _(){}
		func _(){}
		var _, _ int`, "33")

	// Bugs found when writing OS8
	o(` func main() { a, b := f(); printInt(a); printInt(len(b)) }
		var dat [5]int
		func f() (int, []int) { return 33, dat[:] }`, "33\n5")
	o(` func p(a []uint) {}
		func main() {
			var t [10]uint; t2:=t[:];
			before := uint(&t2[0]); p(t2[2:5]); after := uint(&t2[0])
			if before != after { panic() }
		}`, "")
}

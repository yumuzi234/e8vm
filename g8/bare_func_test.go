package g8

import (
	"testing"

	"strings"

	"e8vm.io/e8vm/arch8"
)

func bareTestRun(t *testing.T, input string, N int) (out string, e error) {
	bs, es, _ := CompileBareFunc("main.g", input)
	if es != nil {
		t.Log(input)
		for _, e := range es {
			t.Log(e)
		}
		t.Error("compile failed")
		return "", errRunFailed
	}

	ncycle, out, e := arch8.RunImageOutput(bs, N)
	if ncycle == N {
		t.Log(input)
		t.Error("running out of time")
		return "", errRunFailed
	}
	return out, e
}

func TestBareFunc_good(t *testing.T) {
	const N = 100000
	o := func(input, output string) {
		out, e := bareTestRun(t, input, N)
		if e == errRunFailed {
			t.Error(e)
			return
		}
		if !arch8.IsHalt(e) {
			t.Log(input)
			t.Log(e)
			t.Error("did not halt gracefully")
			return
		}
		out = strings.TrimSpace(out)
		output = strings.TrimSpace(output)
		if out != output {
			t.Log(input)
			t.Logf("expect: %s", output)
			t.Errorf("got: %s", out)
		}
	}

	o(";;;;", "")
	o("printInt(0)", "0")
	o("printInt(3)", "3")
	o("printInt(-444)", "-444")
	o("printInt(2147483647)", "2147483647")
	o("printInt(-2147483647-1)", "-2147483648")
	o("printInt(-2147483648)", "-2147483648")
	o("printInt(300000000)", "300000000")
	o("printInt(4*5+3)", "23")
	o("printInt(3+4*5)", "23")
	o("printInt((3+4)*5)", "35")
	o("printInt((5*(3+4)))", "35")
	o("printInt(3^1)", "2")
	o("printInt(0xf)", "15")
	o("printInt(0xA)", "10")
	o("a:=3; if a==3 { printInt(5) }", "5")
	o("a:=5; if a==3 { printInt(5) }", "")
	o("a:=5; if a==3 { printInt(5) } else { printInt(10) }", "10")
	o("a:=3; for a>0 { printInt(a); a-- }", "3\n2\n1")
	o("a:=0; for a<4 { printInt(a); a++ }", "0\n1\n2\n3")
	o("for i:=0;i<3;i++ { printInt(i); }", "0\n1\n2")
	o("for i:=0;i<10;i+=3 { printInt(i); }", "0\n3\n6\n9")
	o("i:=3; for i:=0;i<3;i++ { printInt(i); }", "0\n1\n2")
	o("i:=0; for ;i<3;i++ { printInt(i); }", "0\n1\n2")
	o("a:=1; { a:=3; printInt(a) }", "3")
	o("true:=3; printInt(true)", "3")
	o("a,b:=3,4; printInt(a); printInt(b)", "3\n4")
	o("a,b:=3,4; { a,b:=b,a; printInt(a); printInt(b) }", "4\n3")
	o("a,b:=3,4; a,b=b,a; printInt(a); printInt(b)", "4\n3")
	o("var a int; printInt(a)", "0")
	o("var a (int); printInt(a)", "0")
	o("var a,b = 3,4; printInt(a); printInt(b)", "3\n4")
	o("var a,b = 3,4; printInt(a); printInt(b)", "3\n4")
	o("var a,b int = 3,4; printInt(a); printInt(b)", "3\n4")
	o("var a,b uint = 3,4; printUint(a); printUint(b)", "3\n4")
	o(` a,b:=3,4; { var a,b=b,a; printInt(a); printInt(b) }
	   	printInt(a); printInt(b)
	`, "4\n3\n3\n4")
	o("var i int; for i < 3 { printInt(i); i=i+1 }", "0\n1\n2")
	o("for true { break }; printInt(3)", "3")
	o("for true { if true break }; printInt(3)", "3")
	o("for { break }; printInt(33)", "33")
	o("for i:=0; i<5; i++ { printInt(i); i++; continue }", "0\n2\n4")
	o("i:=0; for i<3 { printInt(i); i=i+1; continue; break }", "0\n1\n2")
	o("printChar('x')", "x")
	o("var a=32; var b=*&a; printInt(b)", "32")
	o("var a=32; var b=&a; var c=*b; printInt(c)", "32")
	o("var a=32; var b int = *&*&a; printInt(b)", "32")
	o("var a='x'; var b = *&*&a; printChar(b)", "x")
	o("if nil==nil { printInt(3) }", "3")
	o("if nil!=nil { printInt(3) }", "")
	o("var a*int; if a==nil { printInt(3) }", "3")
	o("var a*int; if nil==a { printInt(3) }", "3")
	o("b:=3; a:=&b; if *a==3 { printInt(*a) }", "3")
	o("b:=3; a:=&b; a=nil; if a==nil { printInt(b) }", "3")
	o("b:=3; a:=&b; *a=4; printInt(b)", "4")
	o("if true==true { printChar('y') }", "y")
	o("if true!=true { printChar('y') }", "")
	o("if true==false { printChar('y') }", "")
	o("if true!=false { printChar('y') }", "y")
	o("if false==false { printChar('y') }", "y")
	o("if false!=false { printChar('y') }", "")
	o("var a [4]int; a[3] = 3; printInt(a[3]); printInt(a[2])", "3\n0")
	o("var a [7]int; a[3]=33; pt:=&a[3]; printInt(*pt)", "33")
	o("var a [7]int; printInt(len(a))", "7")
	o("var a [7]int; s:=a[:]; printInt(len(s))", "7")
	o("var a [7]int; s:=a[:3]; printInt(len(s))", "3")
	o("var a [7]int; s:=a[1:]; printInt(len(s))", "6")
	o("var a [7]int; s:=a[1:3]; printInt(len(s))", "2")
	o("var a [7]int; s:=a[0:0]; printInt(len(s))", "0")
	o("var a [7]int; s:=a[:]; a[3]=33; printInt(s[3])", "33")
	o("var a [7]int; s:=a[1:]; a[3]=33; printInt(s[2])", "33")
	o("var a [7]int; s:=a[1:4]; a[3]=33; printInt(s[2])", "33")
	o("var a [7]int; s:=a[:]; a[3]=33; pt:=&s[3]; printInt(*pt)", "33")
	o("a:=3; a++; printInt(a)", "4")
	o("a:=3; pt := &a; *pt++; printInt(a)", "4")

	o("printInt(int(byte(int(-1))))", "255")
	o("printInt(int(byte(3)))", "3")
	o("printInt(int(byte(int(256))))", "0")
	o("printInt(int(char(int(-1))))", "-1")
	o("printInt(int(int8(-1)))", "-1")
	o("printInt(int(char(int(255))))", "-1")
	o("printInt(int(char(int(256))))", "0")
	o("printInt(int(byte(char(int(255)))))", "255")
	o("printInt(int(byte(char(int(-1)))))", "255")

	o("var a int8=-3; printInt(int(a))", "-3")
	o("var a=[]int8{-3}; printInt(int(a[0]))", "-3")

	o("printInt(33 << uint(1))", "66")
	o("printInt(33 << 1)", "66")
	o("printInt(33 >> uint(1))", "16")
	o("printInt(33 >> 1)", "16")
	o("printInt(-33 >> uint(1))", "-17")
	o("printInt(-1 >> uint(1))", "-1")
	o("printInt(int(byte(255) << uint(1)))", "254")
	o("printInt(int(byte(255) << 1))", "254")
	o("printInt(int(uint(33) >> uint(1)))", "16")
	o("printInt(int(uint(33) >> 1))", "16")
	o("printUint(uint(0x80000000) / 10)", "214748364")
	o("printUint(uint(0x80000000) % 10)", "8")
	o("a:=uint(214748364); printUint(a*2)", "429496728")

	o("a:=3; a+=4; printInt(a)", "7")
	o("a:=3; a-=4; printInt(a)", "-1")
	o("a:=3; a*=4; printInt(a)", "12")
	o("a:=3; a/=2; printInt(a)", "1")
	o("a:=uint(3); a/=2; printUint(a)", "1")
	o("a:=33; a<<=uint(1); printInt(a)", "66")
	o("a:=33; a<<=1; printInt(a)", "66")
	o("a:=33; a>>=uint(1); printInt(a)", "16")
	o("a:=33; a>>=1; printInt(a)", "16")

	o("a:=33; b:=(*uint)(&a); printUint(*b)", "33")

	o("if 0x33 == 51 { printInt(33) }", "33")

	o("const a = 33; printInt(a)", "33")
	o("const a = 33; printUint(a)", "33")
	o("const ( a,b=3,4; c=a+b ); printInt(a+b+c)", "14")
	o("const a,b=3,4; var v [a+b]int; printInt(len(v))", "7")

	o("var a []int = nil; printInt(len(a))", "0")
	o("a := []int{}; printInt(len(a))", "0")
	o("var a byte; if a == 0 { printInt(33) }", "33")
}

func TestBareFunc_bad(t *testing.T) {
	// compile errors

	o := func(input string) {
		_, es, _ := CompileBareFunc("main.g", input)
		if es == nil {
			t.Log(input)
			t.Error("should error")
			return
		}
	}

	o("a")                   // expression statement
	o("printInt")            // expression statement
	o("3+4")                 // expression statement
	o("a=3")                 // a not defined
	o("3=4")                 // non-addressable
	o("3=a")                 // non-addressable
	o("var x = &3")          // non-addressable
	o("var a, b int; a+b=3") // non-addressable
	o("a:=3;a:=4")           // redefine
	o("printInt(true)")      // type mismatch
	o("printInt(3, 4)")      // arg count mismatch
	o("printInt()")          // arg count mismatch
	o("a := printInt(3, 4)") // mismatch
	o("a := 3, 4")           // count mismatch
	o("a, b := 3")           // count mismatch
	o("a, b := ()")          // invalid
	o("a()")                 // undefined function
	o("var a, b c")          // undefined type
	o("var a int; var b a")  // not a type
	o("var a = nil")         // infer type from nil
	o("a := nil")            // inter type from nil
	o("break")               // not in for loop
	o("continue")            // not in for loop
	o("if true { break }")   // nothing to break
	o("if true break")       // nothing to break
	o("true > false")        // boolean cannot compare
	o("true + 3")            // boolean cannot add
	o("3++")                 // inc on non-addressable
	o("(3)+=3")              // assign to non-addressable
	o("a := int")

	o("var a [8]int; i:=a[-1]") // negative array index
	o("var a [7]int; s:=a[:]; i:=s[-33]")
	o("var a [0==0]int")

	o("var a int; var b = &a+3")      // pointer cannot add
	o("var a int; var b *uint = &a;") // incompatible pointer type

	o("a := 3/0") // divide by zero
	o("a := 3%0")

	o("const a = -1; printUint(a)")
}

func TestBareFunc_panic(t *testing.T) {
	// runtime errors

	const N = 100000
	o := func(input string) {
		_, e := bareTestRun(t, input, N)
		if !arch8.IsErr(e, arch8.ErrPanic) {
			t.Log(input)
			t.Log(e)
			t.Error("should panic")
			return
		}
	}

	o("var a *int; var b=*a")
	o("var a *int=nil; var b=*a")
	o("var a *int; printInt(*a)")
	o("var a *bool; if *a {}")
	o("var s []int; i:=s[0]")
	o("var a [8]int; j:=-1; i:=a[j]")
	o("var a [8]int; i:=a[9]")
	o("var a [7]int; s:=a[0:0]; i:=s[0]")
	o("var a [7]int; s:=a[1:3]; i:=s[2]")
	o("var a [7]int; s:=a[:]; j:=-33; i:=s[j]")
	o("d:=0; a:=3/d")
	o("d:=0; a:=-3%d")
	o("var d [3]int; s:=d[:]; s=nil; printInt(s[1])")
	o("d:=make([]int, 3, uint(0))")
}

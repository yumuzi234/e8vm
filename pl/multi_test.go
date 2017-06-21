package pl

import (
	"strings"
	"testing"

	"shanhu.io/smlvm/arch"
)

func multiTestRun(t *testing.T, fs map[string]string, N int) (
	string, error,
) {
	bs, es, _ := buildMulti(false, fs, nil)
	if es != nil {
		t.Log(fs)
		for _, err := range es {
			t.Log(err)
		}
		t.Error("compile failed")
		return "", errRunFailed
	}

	ncycle, out, err := arch.RunImageOutput(bs, N)
	if ncycle == N {
		t.Log(fs)
		t.Error("running out of time")
		return "", errRunFailed
	}
	return out, err
}

func TestMultiFile(t *testing.T) {
	const N = 100000

	o := func(fs map[string]string, output string) {
		out, err := multiTestRun(t, fs, N)
		if err == errRunFailed {
			t.Error(err)
			return
		}

		if !arch.IsHalt(err) {
			t.Log(fs)
			t.Log(err)
			t.Error("did not halt gracefully")
			return
		}

		got := strings.TrimSpace(out)
		expect := strings.TrimSpace(output)
		if got != expect {
			t.Log(fs)
			t.Logf("expect: %s", expect)
			t.Errorf("got: %s", got)
		}
	}
	type files map[string]string

	o(files{
		"main/m.g": "func main() { }",
	}, "")

	o(files{
		"a/a.g": "func P() { printInt(33) }",
		"main/m.g": `
			import ( "a" )
			func main() { a.P() }`,
	}, "33")

	o(files{
		"a/a.g":    "struct A {}; func (a *A) P() { printInt(33) }",
		"b/b.g":    `import ("a"); var A a.A`,
		"main/m.g": `import ("b"); func main() { b.A.P() }`,
	}, "33")

	o(files{
		"a/a.g":    "func init() { printInt(33) }",
		"b/b.g":    `import (_ "a"); func init() { printInt(44) }`,
		"main/m.g": `import (_ "b"); func main() { printInt(55) }`,
	}, "33\n44\n55")

	o(files{
		"a/a.g": "const A=33",
		"main/m.g": `
			import ("a")
			var array [a.A]int
			func main() { printInt(len(array)) }`,
	}, "33")
	o(files{
		"a/a.g": "const A=33+5-2",
		"main/m.g": `
			import ("a")
			var array [a.A-3]int
			func main() { printInt(len(array)) }`,
	}, "33")

	o(files{
		"asm/a/a.s": `
			func F {
				mov pc ret
			}`,
		"main/m.g": `
			import ("asm/a")
			func main() { a.F(); printInt(33) }`,
	}, "33")

	o(files{
		"asm/a/a.s": `
			func F {
				addi r1 r0 33
				mov pc ret
			}`,
		"main/m.g": `
			import ("asm/a")
			func f() int = a.F
			func main() { printInt(f()) }`,
	}, "33")

	// A bug found when writing mempair.
	o(files{
		"a/a.g": `struct A { I int }`,
		"main/m.g": `
			import ("a")
			func f() (int, a.A) {
				var ret a.A
				ret.I = 33
				return 0, ret
			}
			func main() {
				_, v := f();
				printInt(v.I)
			}`,
	}, "33")
}

func TestMultiFileBad(t *testing.T) {
	o := func(files map[string]string) {
		_, es, _ := buildMulti(false, files, nil)
		if es == nil {
			t.Error("should error")
			return
		}

		t.Log(files)
		for _, e := range es {
			t.Log(e)
		}
	}
	type files map[string]string

	// circular dependency
	o(files{
		"main/a.g": `func main() { a() }; func b() {}`,
		"main/b.g": `func a() { b() }`,
	})

	// name conflict
	o(files{
		"a/a.g":    ``,
		"main/a.g": `import ("a"); func a() {}; func main() {};`,
	})

	// unused import
	o(files{
		"a/a.g":    ``,
		"main/a.g": `import ("a"); func main() {};`,
	})

	// using private methods
	o(files{
		"a/a.g":    `func f() {}`,
		"main/a.g": `import ("a"); func main() { a.f() };`,
	})

	o(files{
		"asm/a/a.g": `
			func A {
			}
		`,
		"main/a.g": `
			import ("asm/a")
			struct A { func f() = a.A; }
			func main() { }
		`,
	})
}

package g8

import (
	"path"
	"strings"
	"testing"

	"e8vm.io/e8vm/arch8"
	"e8vm.io/e8vm/build8"
	"e8vm.io/e8vm/lex8"
)

func buildMulti(files map[string]string) (
	image []byte, errs []*lex8.Error, log []byte,
) {
	home := makeMemHome(Lang())

	pkgs := make(map[string]*build8.MemPkg)
	for f, content := range files {
		p := path.Dir(f)
		base := path.Base(f)
		pkg, found := pkgs[p]
		if !found {
			pkg = home.NewPkg(p)
		}
		pkg.AddFile(f, base, content)
	}

	return buildMainPkg(home)
}

func multiTestRun(t *testing.T, fs map[string]string, N int) (
	string, error,
) {
	bs, es, _ := buildMulti(fs)
	if es != nil {
		t.Log(fs)
		for _, err := range es {
			t.Log(err)
		}
		t.Error("compile failed")
		return "", errRunFailed
	}

	ncycle, out, err := arch8.RunImageOutput(bs, N)
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

		if !arch8.IsHalt(err) {
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
		"main/main.g": "func main() { }",
	}, "")

	o(files{
		"a/a.g": "func P() { printInt(33) }",
		"main/main.g": `
			import ( "a" )
			func main() { a.P() }`,
	}, "33")

	o(files{
		"a/a.g":       "struct A { func P() { printInt(33) } }",
		"b/b.g":       `import ("a"); var A a.A`,
		"main/main.g": `import ("b"); func main() { b.A.P() }`,
	}, "33")

	o(files{
		"a/a.g":       "func init() { printInt(33) }",
		"b/b.g":       `import (_ "a"); func init() { printInt(44) }`,
		"main/main.g": `import (_ "b"); func main() { printInt(55) }`,
	}, "33\n44\n55")
}

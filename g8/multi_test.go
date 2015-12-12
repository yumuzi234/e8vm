package g8

import (
	"path"
	"testing"

	"e8vm.io/e8vm/build8"
	"e8vm.io/e8vm/arch8"
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

func testMulti(t *testing.T, fs map[string]string, N int) (string, error) {
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

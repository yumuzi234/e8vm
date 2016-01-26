package g8

import (
	"path"

	"e8vm.io/e8vm/build8"
	"e8vm.io/e8vm/lex8"
)

func buildMulti(lang build8.Lang, files map[string]string, runTests bool) (
	image []byte, errs []*lex8.Error, log []byte,
) {
	home := makeMemHome(lang)

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

	return buildMainPkg(home, runTests)
}

// CompileMulti compiles a set of files into a bare-metal E8 image
func CompileMulti(files map[string]string, golike, runTests bool) (
	[]byte, []*lex8.Error,
) {
	var lang build8.Lang
	if !golike {
		lang = Lang()
	} else {
		lang = LangGoLike()
	}

	image, errs, _ := buildMulti(lang, files, runTests)
	return image, errs
}

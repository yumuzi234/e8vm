package g8

import (
	"errors"
	"path/filepath"

	"e8vm.io/e8vm/asm8"
	"e8vm.io/e8vm/build8"
	"e8vm.io/e8vm/lex8"
)

func buildSingle(fname, s string, lang build8.Lang) (
	image []byte, es []*lex8.Error, log []byte,
) {
	home := build8.NewMemHome(lang)
	home.AddLang("asm", asm8.Lang())

	pkg := home.NewPkg("main")
	name := filepath.Base(fname)
	pkg.AddFile(fname, name, s)

	builtin := home.NewPkg("asm/builtin")
	builtin.AddFile("", "builtin.s", builtInSrc)

	b := build8.NewBuilder(home)
	es = b.BuildAll()
	if es != nil {
		return nil, es, nil
	}

	image = home.Bin("main")
	log = home.Log("main", "ir")
	if image == nil {
		err := errors.New("missing main() function, no binary created")
		return nil, lex8.SingleErr(err), log
	}

	return image, nil, log
}

// CompileSingle compiles a file into a bare-metal E8 image
func CompileSingle(fname, s string, golike bool) (
	[]byte, []*lex8.Error, []byte,
) {
	var lang build8.Lang
	if !golike {
		lang = Lang()
	} else {
		lang = LangGolike()
	}
	return buildSingle(fname, s, lang)
}

package g8

import (
	"e8vm.io/e8vm/build8"
	"e8vm.io/e8vm/lex8"
)

func buildMulti(
	lang build8.Lang, files map[string]string, opt *build8.Options,
) (
	image []byte, errs []*lex8.Error, log []byte,
) {
	home := MakeMemHome(lang)
	home.AddFiles(files)
	return buildMainPkg(home, opt)
}

// CompileMulti compiles a set of files into a bare-metal E8 image
func CompileMulti(
	files map[string]string, golike bool, opt *build8.Options,
) (
	[]byte, []*lex8.Error,
) {
	lang := Lang(golike)
	image, errs, _ := buildMulti(lang, files, opt)
	return image, errs
}

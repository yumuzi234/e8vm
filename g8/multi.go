package g8

import (
	"e8vm.io/e8vm/build8"
	"e8vm.io/e8vm/lex8"
)

func buildMulti(lang build8.Lang, files map[string]string, runTests bool) (
	image []byte, errs []*lex8.Error, log []byte,
) {
	home := MakeMemHome(lang)
	home.AddFiles(files)
	return buildMainPkg(home, runTests, 0)
}

// CompileMulti compiles a set of files into a bare-metal E8 image
func CompileMulti(files map[string]string, golike, runTests bool) (
	[]byte, []*lex8.Error,
) {
	lang := Lang(golike)
	image, errs, _ := buildMulti(lang, files, runTests)
	return image, errs
}

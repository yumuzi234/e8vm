package pl

import (
	"shanhu.io/smlvm/builds"
	"shanhu.io/smlvm/lexing"
)

func buildMulti(
	lang *builds.Lang, files map[string]string, opt *builds.Options,
) (
	image []byte, errs []*lexing.Error, log []byte,
) {
	home := MakeMemHome(lang)
	home.AddFiles(files)
	return buildMainPkg(home, opt)
}

// CompileMulti compiles a set of files into a bare-metal E8 image
func CompileMulti(
	files map[string]string, golike bool, opt *builds.Options,
) (
	[]byte, []*lexing.Error,
) {
	lang := Lang(golike)
	image, errs, _ := buildMulti(lang, files, opt)
	return image, errs
}

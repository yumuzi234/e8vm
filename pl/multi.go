package pl

import (
	"shanhu.io/smlvm/builds"
	"shanhu.io/smlvm/lexing"
)

func buildMulti(golike bool, files map[string]string, opt *builds.Options) (
	image []byte, errs []*lexing.Error, log []byte,
) {
	fs := MakeMemFS()
	lp := MakeLangSet(golike)
	for f, s := range files {
		fs.AddTextFile(f, s)
	}
	return buildMainPkg(fs, lp, opt)
}

// CompileMulti compiles a set of files into a bare-metal E8 image
func CompileMulti(
	files map[string]string, golike bool, opt *builds.Options,
) (
	[]byte, []*lexing.Error,
) {
	image, errs, _ := buildMulti(golike, files, opt)
	return image, errs
}

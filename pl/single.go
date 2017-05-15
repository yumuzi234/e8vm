package pl

import (
	"path"
	"path/filepath"
	"strings"

	"shanhu.io/smlvm/asm"
	"shanhu.io/smlvm/builds"
	"shanhu.io/smlvm/lexing"
)

// MakeMemFS makes a memory filesystem for compiling.
// It contains the basic built-in packages.
func MakeMemFS() *builds.MemFS {
	home := builds.NewMemFS()
	builtInDir := strings.TrimPrefix(BuiltInPkg, "/")

	if err := home.MakeDir(builtInDir); err != nil {
		panic(err)
	}

	err := home.AddTextFile(path.Join(builtInDir, "builtin.s"), BuiltInSrc)
	if err != nil {
		panic(err)
	}

	return home
}

// MakeLangPicker makes the language picker using the given language as the
// default language and assembly for "asm" keyword.
func MakeLangPicker(lang *builds.Lang) *builds.LangPicker {
	ret := builds.NewLangPicker(lang)
	ret.AddLang("asm", asm.Lang())
	return ret
}

func buildMainPkg(
	fs *builds.MemFS, langPicker *builds.LangPicker,
	opt *builds.Options,
) (image []byte, errs []*lexing.Error, log []byte) {
	out := builds.NewMemFS()
	b := builds.NewBuilder(fs, langPicker, "", out)
	if opt != nil {
		b.Options = opt
	}
	if errs := b.BuildAll(); errs != nil {
		return nil, errs, nil
	}

	ok, err := out.HasFile("bin/main.bin")
	if err != nil {
		return nil, lexing.SingleErr(err), nil
	}
	if !ok {
		return nil, lexing.SingleCodeErr("pl.missingMainFunc", err), nil
	}

	image, err = out.Read("bin/main.bin")
	if err != nil {
		return nil, lexing.SingleErr(err), nil
	}
	log, err = out.Read("out/main/ir")
	if err != nil {
		return nil, lexing.SingleErr(err), nil
	}

	return image, nil, log
}

func buildSingle(
	f, s string, lang *builds.Lang, opt *builds.Options,
) (
	image []byte, errs []*lexing.Error, log []byte,
) {
	fs := MakeMemFS()
	err := fs.AddTextFile(path.Join("main", filepath.Base(f)), s)
	if err != nil {
		return nil, lexing.SingleErr(err), nil
	}
	lp := MakeLangPicker(lang)
	return buildMainPkg(fs, lp, opt)
}

// CompileSingle compiles a file into a bare-metal E8 image
func CompileSingle(fname, s string, golike bool) (
	[]byte, []*lexing.Error, []byte,
) {
	return buildSingle(fname, s, Lang(golike), new(builds.Options))
}

// CompileAndTestSingle compiles a file into a bare-metal E8 image and
// runs the tests.
func CompileAndTestSingle(fname, s string, golike bool, testCycles int) (
	[]byte, []*lexing.Error, []byte,
) {
	return buildSingle(fname, s, Lang(golike), &builds.Options{
		RunTests:   true,
		TestCycles: testCycles,
	})
}

package pl

import (
	"errors"
	"path/filepath"

	"e8vm.io/e8vm/asm"
	"e8vm.io/e8vm/builds"
	"e8vm.io/e8vm/lexing"
)

// MakeMemHome makes a memory home for compiling.
// It contains the basic built-in packages.
func MakeMemHome(lang builds.Lang) *builds.MemHome {
	home := builds.NewMemHome(lang)
	home.AddLang("asm", asm.Lang())
	builtin := home.NewPkg("asm/builtin")
	builtin.AddFile("", "builtin.s", BuiltInSrc)

	return home
}

func buildMainPkg(home *builds.MemHome, opt *builds.Options) (
	image []byte, errs []*lexing.Error, log []byte,
) {
	b := builds.NewBuilder(home, home)
	if opt != nil {
		b.Options = opt
	}
	if errs := b.BuildAll(); errs != nil {
		return nil, errs, nil
	}

	image = home.BinBytes("main")
	log = home.OutputBytes("main", "ir")
	if image == nil {
		err := errors.New("missing main() function, no binary created")
		return nil, lexing.SingleErr(err), log
	}

	return image, nil, log
}

func buildSingle(
	f, s string, lang builds.Lang, opt *builds.Options,
) (
	image []byte, errs []*lexing.Error, log []byte,
) {
	home := MakeMemHome(lang)

	pkg := home.NewPkg("main")
	name := filepath.Base(f)
	pkg.AddFile(f, name, s)

	return buildMainPkg(home, opt)
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

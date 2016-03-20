package g8

import (
	"errors"
	"path/filepath"

	"e8vm.io/e8vm/asm8"
	"e8vm.io/e8vm/build8"
	"e8vm.io/e8vm/lex8"
)

// MakeMemHome makes a memory home for compiling.
// It contains the basic built-in packages.
func MakeMemHome(lang build8.Lang) *build8.MemHome {
	home := build8.NewMemHome(lang)
	home.AddLang("asm", asm8.Lang())
	builtin := home.NewPkg("asm/builtin")
	builtin.AddFile("", "builtin.s", builtInSrc)

	return home
}

func buildMainPkg(home *build8.MemHome, opt *build8.Options) (
	image []byte, errs []*lex8.Error, log []byte,
) {
	b := build8.NewBuilder(home, home)
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
		return nil, lex8.SingleErr(err), log
	}

	return image, nil, log
}

func buildSingle(
	f, s string, lang build8.Lang, opt *build8.Options,
) (
	image []byte, errs []*lex8.Error, log []byte,
) {
	home := MakeMemHome(lang)

	pkg := home.NewPkg("main")
	name := filepath.Base(f)
	pkg.AddFile(f, name, s)

	return buildMainPkg(home, opt)
}

// CompileSingle compiles a file into a bare-metal E8 image
func CompileSingle(fname, s string, golike bool) (
	[]byte, []*lex8.Error, []byte,
) {
	return buildSingle(fname, s, Lang(golike), new(build8.Options))
}

// CompileAndTestSingle compiles a file into a bare-metal E8 image and
// runs the tests.
func CompileAndTestSingle(fname, s string, golike bool, testCycles int) (
	[]byte, []*lex8.Error, []byte,
) {
	return buildSingle(fname, s, Lang(golike), &build8.Options{
		RunTests:   true,
		TestCycles: testCycles,
	})
}

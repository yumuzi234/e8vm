package build8

import (
	"e8vm.io/e8vm/dagvis"
	"e8vm.io/e8vm/lex8"
)

// Options contains options for building a set of packages.
type Options struct {
	Verbose bool
	InitPC  uint32

	StaticOnly bool
	RunTests   bool
	TestCycles int

	SaveDeps       func(deps *dagvis.Map)
	SaveFileTokens func(p string, toks []*lex8.Token)
	LogLine        func(s string)
}

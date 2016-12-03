package builds

import (
	"shanhu.io/smlvm/dagvis"
	"shanhu.io/smlvm/lexing"
)

// Options contains options for building a set of packages.
type Options struct {
	Verbose bool
	InitPC  uint32

	StaticOnly bool
	RunTests   bool
	TestCycles int

	SaveDeps       func(deps *dagvis.Graph)
	SaveFileTokens func(p string, toks []*lexing.Token)
	LogLine        func(s string)
}

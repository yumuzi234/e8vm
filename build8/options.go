package build8

import (
	"e8vm.io/e8vm/dagvis"
)

// Options contains options for building a set of packages.
type Options struct {
	Verbose bool
	InitPC  uint32

	StaticOnly bool
	RunTests   bool

	SaveDeps func(deps *dagvis.Map)
}

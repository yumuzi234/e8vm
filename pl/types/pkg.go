package types

import (
	"fmt"

	"e8vm.io/e8vm/syms"
)

// Pkg represents a package import.
type Pkg struct {
	As   string
	Lang string
	Syms *syms.Table
}

// Size will panic.
func (p *Pkg) Size() int32 { panic("bug") }

// RegSizeAlign will panic.
func (p *Pkg) RegSizeAlign() bool { panic("bug") }

func (p *Pkg) String() string { return fmt.Sprintf("package %s", p.As) }

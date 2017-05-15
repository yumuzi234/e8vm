package arch

import (
	"io"
	"time"

	"shanhu.io/smlvm/arch/devs"
	"shanhu.io/smlvm/net"
)

// Config contains config for constructing a machine
type Config struct {
	MemSize uint32
	Ncore   int

	Output   io.Writer
	Net      net.Handler
	Screen   devs.ScreenRender
	RandSeed int64

	InitPC       uint32
	InitSP       uint32
	StackPerCore uint32

	BootArg uint32

	ROM string

	// time functions
	Now     func() time.Time
	PerfNow func() time.Duration
}

package arch8

import (
	"io"
)

// Config contains config for constructing a machine
type Config struct {
	MemSize uint32
	Ncore   int

	Output   io.Writer
	Monitor  Monitor
	RandSeed int64

	InitPC       uint32
	InitSP       uint32
	StackPerCore uint32

	BootArg uint32

	ROM string
}

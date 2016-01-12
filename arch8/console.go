package arch8

import (
	"io"
	"log"
	"os"
)

// Console is a simple console that can output/input a single
// byte at a time
type console struct {
	intBus intBus
	p      *pageOffset

	Core   byte
	IntIn  byte
	IntOut byte

	Output io.Writer
}

// NewConsole creates a new simple console.
func newConsole(p *page, i intBus) *console {
	ret := new(console)
	ret.intBus = i
	const consoleBase = 0
	ret.p = &pageOffset{p, consoleBase}

	ret.Core = 0
	ret.IntIn = 8
	ret.IntOut = 9

	ret.Output = os.Stdout
	return ret
}

const (
	consoleOut      = 0
	consoleOutValid = 1

	consoleIn      = 4
	consoleInValid = 5
)

func (c *console) interrupt(code byte) {
	c.intBus.Interrupt(code, c.Core)
}

// Tick flushes the buffered byte to the console.
func (c *console) Tick() {
	outValid := c.p.readByte(consoleOutValid)
	if outValid != 0 {
		out := c.p.readByte(consoleOut)
		_, e := c.Output.Write([]byte{out})
		if e != nil {
			log.Print(e)
		}
		c.p.writeByte(consoleOutValid, 0)
		c.interrupt(c.IntOut) // out available
	}

	// TODO(h8liu): input part
}

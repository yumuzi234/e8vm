package arch

import (
	"io"
	"log"
	"os"

	"shanhu.io/smlvm/arch/devs"
)

// Console is a simple console that can output/input a single
// byte at a time
type console struct {
	intBus intBus
	p      *pageOffset

	Core   byte // Core to throw exception
	IntIn  byte
	IntOut byte

	Output io.Writer
}

// NewConsole creates a new simple console.
func newConsole(p *page, i intBus) *console {
	ret := new(console)
	ret.intBus = i
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

func (c *console) Handle(req []byte) ([]byte, int32) {
	const maxOutputLen = 128

	m := len(req)
	if m == 0 {
		return nil, 0
	}
	if m > maxOutputLen {
		return nil, devs.ErrInvalidArg
	}

	if _, err := c.Output.Write(req); err != nil {
		log.Print(err)
	}
	return nil, 0
}

// Tick flushes the buffered byte to the console.
func (c *console) Tick() {
	outValid := c.p.readU8(consoleOutValid)
	if outValid != 0 {
		out := c.p.readU8(consoleOut)
		_, e := c.Output.Write([]byte{out})
		if e != nil {
			log.Print(e)
		}
		c.p.writeU8(consoleOutValid, 0)
		c.interrupt(c.IntOut) // out available
	}

	// TODO(h8liu): input part
}

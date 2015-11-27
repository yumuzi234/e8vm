package arch8

// interrupt defines the interrupt page
type interrupt struct {
	*pageOffset // the dma page for interrupt handler
}

// Ninterrupt is the number of interrupts.
const Ninterrupt = 256

const (
	intFlags     = 0  // flags, bit 0 is master enabling switch
	intHandlerSP = 4  // position of the handler stack base pointer
	intHandlerPC = 8  // position of the handler start PC
	intSyscallSP = 12 // position of the syscall stack base pointer
	intSyscallPC = 16 // position of the syscall start PC
	intMask      = 32 // interrupt enable mask bits offset (32 bytes)
	intPending   = 64 // interrupt pending bits offset (32 bytes)

	intCtrlSize = 128
)

// newInterrupt creates a interrupt on the given DMA page.
func newInterrupt(p *page, core byte) *interrupt {
	ret := new(interrupt)
	base := uint32(core) * intCtrlSize
	if base+intCtrlSize > PageSize {
		panic("bug")
	}
	ret.pageOffset = &pageOffset{p, base}

	return ret
}

func (in *interrupt) handlerSP() uint32 { return in.readWord(intHandlerSP) }
func (in *interrupt) handlerPC() uint32 { return in.readWord(intHandlerPC) }
func (in *interrupt) syscallSP() uint32 { return in.readWord(intSyscallSP) }
func (in *interrupt) syscallPC() uint32 { return in.readWord(intSyscallPC) }

// Issue issues an interrupt. If the interrupt is already issued,
// this has no effect.
func (in *interrupt) Issue(i byte) {
	off := uint32(i/8) + intPending
	b := in.readByte(off)
	b |= 0x1 << (i % 8)
	in.writeByte(off, b)
}

// Clear clears an interrupt. If the interrupt is already cleared,
// this has no effect.
func (in *interrupt) Clear(i byte) {
	off := uint32(i/8) + intPending
	b := in.readByte(off)
	b &= ^(0x1 << (i % 8))
	in.writeByte(off, b)
}

// Enable sets the interrupt enable bit in the flags.
func (in *interrupt) Enable() {
	b := in.readByte(intFlags)
	b |= 0x1
	in.writeByte(intFlags, b)
}

// Enabled tests if interrupt is enabled
func (in *interrupt) Enabled() bool {
	b := in.readByte(intFlags)
	return (b & 0x1) != 0
}

// Disable clears the interrupt enable bit in the flags.
func (in *interrupt) Disable() {
	b := in.readByte(intFlags)
	b &= ^byte(0x1)
	in.writeByte(intFlags, b)
}

// EnableInt enables a particular interrupt by clearing the mask.
func (in *interrupt) EnableInt(i byte) {
	off := uint32(i/8) + intMask
	b := in.readByte(off)
	b |= 0x1 << (i % 8)
	in.writeByte(off, b)
}

// DisableInt disables a particular interrupt by setting the mask.
func (in *interrupt) DisableInt(i byte) {
	off := uint32(i/8) + intMask
	b := in.readByte(off)
	b &= ^(0x1 << (i % 8))
	in.writeByte(off, b)
}

// Flags returns the flags byte.
func (in *interrupt) Flags() byte {
	return in.readByte(intFlags)
}

// Poll looks for the next pending interrupt.
func (in *interrupt) Poll() (bool, byte) {
	flag := in.Flags()
	if flag&0x1 == 0 { // interrupt is disabled
		return false, 0
	}

	// search bits based on priorities.
	// smaller is higher
	for i := uint32(0); i < Ninterrupt/32; i++ {
		pending := in.readWord(intPending + i*4)
		mask := in.readWord(intMask + i*4)
		pending &= mask
		if pending == 0 {
			continue
		}

		for b := byte(0); b < 32; b++ {
			if pending&(0x1<<b) == 0 {
				continue
			}

			return true, byte(i*32) + b
		}
	}

	return false, 0
}

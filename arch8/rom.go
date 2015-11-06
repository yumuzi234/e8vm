package arch8

const (
	romCmd     = 0
	romNameLen = 1
	romState   = 2
	romErr     = 3

	romOffset = 4
	romAddr   = 8
	romSize   = 12
	romNread  = 16

	romFilename    = 20
	romFilenameMax = 100
)

const (
	romCmdIdle    = 0
	romCmdRequest = 1

	romStateIdle = 0
	romStateBusy = 1

	romErrNone     = 0
	romErrNotFound = 1
	romErrEOF      = 2
)

type rom struct {
	intBus intBus
	p      *pageOffset
	mem    *phyMemory
	root   string

	state     byte
	countDown int

	Core byte
}

func newRom(p *page, mem *phyMemory, i intBus, root string) *rom {
	return &rom{
		intBus: i,
		p:      &pageOffset{p, 0x100},
		mem:    mem,
		root:   root,
	}
}

func (r *rom) interrupt(code byte) {
	r.intBus.Interrupt(code, r.Core)
}

func (r *rom) readFile() {
	// TODO:
}

func (r *rom) Tick() {
	switch r.state {
	case romStateIdle:
		cmd := r.p.readByte(romCmd)
		if cmd != 0 {
			r.state = romStateBusy
			r.countDown = 100
		}
	case romStateBusy:
		if r.countDown > 0 {
			r.countDown--
		} else {
			r.readFile()
			r.state = romStateIdle
		}
	}

	r.p.writeByte(romState, r.state)
}

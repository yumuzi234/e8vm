package arch

import (
	"io"
	"log"
	"os"
	"path/filepath"
)

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
)

const (
	romErrNone = iota
	romErrEOF
	romErrNotFound
	romErrOpen
	romErrRead
	romErrMemory
)

type rom struct {
	intBus intBus
	p      *pageOffset
	mem    *phyMemory
	root   string

	state byte

	countDown int
	addr      uint32 // bytes to write at
	bs        []byte // bytes read
	err       byte

	Core    byte
	IntDone byte
}

func newROM(p *page, mem *phyMemory, i intBus, root string) *rom {
	return &rom{
		intBus: i,
		p:      &pageOffset{p, romBase},
		mem:    mem,
		root:   root,

		IntDone: IntROM,
	}
}

func (r *rom) interrupt(code byte) {
	r.intBus.Interrupt(code, r.Core)
}

func (r *rom) readFile() (byte, error) {
	nameLen := r.p.readU8(romNameLen)
	offset := r.p.readU32(romOffset)
	addr := r.p.readU32(romAddr)
	size := r.p.readU32(romSize)

	if nameLen > romFilenameMax {
		nameLen = romFilenameMax
	}
	name := make([]byte, nameLen)
	for i := range name {
		name[i] = r.p.readU8(romFilename + uint32(i))
	}

	fullPath := filepath.Join(r.root, string(name))
	f, err := os.Open(fullPath)
	if err != nil {
		if os.IsNotExist(err) {
			return romErrNotFound, err
		}
		return romErrOpen, err
	}
	defer f.Close()

	if _, err = f.Seek(int64(offset), 0); err != nil {
		return romErrRead, err
	}

	buf := make([]byte, size)
	read, err := f.Read(buf)
	if err == io.EOF {
		return romErrEOF, err
	}
	if err != nil {
		return romErrRead, err
	}

	r.addr = addr
	r.bs = buf[:read]

	return 0, nil
}

func (r *rom) Tick() {
	switch r.state {
	case romStateIdle:
		cmd := r.p.readU8(romCmd)
		if cmd != 0 {
			r.state = romStateBusy

			errCode, err := r.readFile()
			if err != nil && err != io.EOF {
				log.Println(err)
			}

			if len(r.bs) == 0 {
				r.countDown = 5
			} else {
				r.countDown = 10 * len(r.bs)
			}

			r.err = errCode
			r.p.writeU8(romCmd, romCmdIdle)
		}
	case romStateBusy:
		if r.countDown > 0 {
			r.countDown--
		} else {
			if r.err == romErrNone {
				r.p.writeU32(romNread, uint32(len(r.bs)))

				for i, b := range r.bs {
					err := r.mem.WriteU8(r.addr+uint32(i), b)
					if err != nil {
						r.err = romErrMemory
						break
					}
				}
			}

			r.p.writeU8(romErr, r.err)
			r.state = romStateIdle
			r.interrupt(r.IntDone)
		}
	}

	r.p.writeU8(romState, r.state)
}

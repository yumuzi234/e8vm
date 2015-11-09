package arch8

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
	romErrMem
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

func (r *rom) readFile() (byte, error) {
	nameLen := r.p.readByte(romNameLen)
	offset := r.p.readWord(romOffset)
	addr := r.p.readWord(romAddr)
	size := r.p.readWord(romSize)

	if nameLen > romFilenameMax {
		nameLen = romFilenameMax
	}
	name := make([]byte, nameLen)
	for i := range name {
		name[i] = r.p.readByte(romFilename + uint32(i))
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
		cmd := r.p.readByte(romCmd)
		if cmd != 0 {
			r.state = romStateBusy
			r.countDown = 100

			errCode, err := r.readFile()
			if err != nil {
				log.Println(err)
			}
			r.err = errCode
		}
	case romStateBusy:
		if r.countDown > 0 {
			r.countDown--
		} else {
			if r.err == romErrNone {
				r.p.writeWord(romNread, uint32(len(r.bs)))

				for i, b := range r.bs {
					err := r.mem.WriteByte(r.addr+uint32(i), b)
					if err != nil {
						r.err = romErrMem
						break
					}
				}
			}

			r.p.writeByte(romErr, r.err)
			r.state = romStateIdle
		}
	}

	r.p.writeByte(romState, r.state)
}

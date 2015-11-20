package arch8

import (
	"bytes"
	"io"
	"math/rand"

	"e8vm.io/e8vm/e8"
)

// Machine is a multicore shared memory simulated arch8 machine.
type Machine struct {
	phyMem  *phyMemory
	inst    inst
	cores   *multiCore
	serial  *serial
	console *console
	ticker  *ticker
	rom     *rom

	devices []device
}

// Default SP settings.
const (
	DefaultSPBase   uint32 = 0x20000
	DefaultSPStride uint32 = 0x2000
)

// NewMachine creates a machine with memory and cores.
// 0 memSize for full 4GB memory.
func NewMachine(memSize uint32, ncore int) *Machine {
	ret := new(Machine)
	ret.phyMem = newPhyMemory(memSize)
	ret.inst = new(instArch8)
	ret.cores = newMultiCore(ncore, ret.phyMem, ret.inst)

	// hook-up devices
	p := ret.phyMem.Page(pageBasicIO)

	ret.serial = newSerial(p, ret.cores)
	ret.console = newConsole(p, ret.cores)
	ret.ticker = newTicker(ret.cores)

	ret.addDevice(ret.ticker)
	ret.addDevice(ret.serial)
	ret.addDevice(ret.console)

	sys := ret.phyMem.Page(pageSysInfo)
	sys.WriteWord(0, ret.phyMem.npage)
	sys.WriteWord(4, uint32(ncore))

	ret.SetSP(DefaultSPBase, DefaultSPStride)

	return ret
}

// MountROM mounts the root of the read-only disk.
func (m *Machine) MountROM(root string) {
	p := m.phyMem.Page(pageBasicIO)
	m.rom = newROM(p, m.phyMem, m.cores, root)
	m.addDevice(m.rom)
}

// WriteByte writes the byte at a particular physical address.
func (m *Machine) WriteByte(phyAddr uint32, b byte) error {
	exp := m.phyMem.WriteByte(phyAddr, b)
	if exp == nil {
		return nil
	}
	return exp
}

// WriteWord writes the word at a particular physical address.
// The address must be word aligned.
func (m *Machine) WriteWord(phyAddr uint32, v uint32) error {
	exp := m.phyMem.WriteWord(phyAddr, v)
	if exp == nil {
		return nil
	}
	return exp
}

// SetOutput sets the output writer of the machine's serial
// console.
func (m *Machine) SetOutput(w io.Writer) {
	m.serial.Output = w
	m.console.Output = w
}

// AddDevice adds a devices to the machine.
func (m *Machine) addDevice(d device) {
	m.devices = append(m.devices, d)
}

// Tick proceeds the simulation by one tick.
func (m *Machine) Tick() *CoreExcep {
	for _, d := range m.devices {
		d.Tick()
	}

	return m.cores.Tick()
}

// Run simulates nticks. It returns the number of ticks
// simulated without error, and the first met error if any.
func (m *Machine) Run(nticks int) (int, *CoreExcep) {
	n := 0
	for i := 0; nticks == 0 || i < nticks; i++ {
		e := m.Tick()
		n++
		if e != nil {
			return n, e
		}
	}

	return n, nil
}

// WriteBytes write a byte buffer to the memory at a particular offset.
func (m *Machine) WriteBytes(r io.Reader, offset uint32) error {
	mod := offset % PageSize
	start := uint32(0)
	if mod != 0 {
		start = PageSize - mod
	}

	pageBuf := make([]byte, PageSize)
	pn := offset / PageSize
	for {
		p := m.phyMem.Page(pn)
		if p == nil {
			return newOutOfRange(offset)
		}

		buf := pageBuf[:PageSize-start]
		n, err := r.Read(buf)
		if err == io.EOF {
			return nil
		}

		p.WriteAt(buf[:n], start)
		start = 0
		pn++
	}

	return nil
}

// RandSeed sets the random seed of the ticker.
func (m *Machine) RandSeed(s int64) {
	m.ticker.Rand = rand.New(rand.NewSource(s))
}

func (m *Machine) loadSections(secs []*e8.Section) error {
	for _, s := range secs {
		var buf io.Reader
		if s.Type == e8.Zeros {
			buf = &zeroReader{s.Size}
		} else {
			buf = bytes.NewReader(s.Bytes)
		}

		if err := m.WriteBytes(buf, s.Addr); err != nil {
			return err
		}
	}
	return nil
}

// SetPC sets all cores to start with a particular PC pointer.
func (m *Machine) SetPC(pc uint32) {
	for _, cpu := range m.cores.cores {
		cpu.regs[PC] = pc
	}
}

// SetSP sets the stack pointer.
func (m *Machine) SetSP(sp, stackSize uint32) {
	for i, cpu := range m.cores.cores {
		cpu.regs[SP] = sp + uint32(i+1)*stackSize
	}
}

func findCodeStart(secs []*e8.Section) (uint32, bool) {
	for _, s := range secs {
		if s.Type == e8.Code {
			return s.Addr, true
		}
	}
	return 0, false
}

// LoadImage loads an e8 image into the machine.
func (m *Machine) LoadImage(r io.ReadSeeker) error {
	secs, err := e8.Read(r)
	if err != nil {
		return err
	}
	if err := m.loadSections(secs); err != nil {
		return err
	}

	if pc, found := findCodeStart(secs); found {
		m.SetPC(pc)
	}
	return nil
}

// LoadImageBytes loads an e8 image in bytes into the machine.
func (m *Machine) LoadImageBytes(bs []byte) error {
	return m.LoadImage(bytes.NewReader(bs))
}

// PrintCoreStatus prints the cpu statuses.
func (m *Machine) PrintCoreStatus() {
	m.cores.PrintStatus()
}

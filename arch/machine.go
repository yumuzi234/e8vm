package arch

import (
	"bytes"
	"fmt"
	"io"
	"math/rand"

	"e8vm.io/e8vm/image"
)

// Machine is a multicore shared memory simulated arch8 machine.
type Machine struct {
	phyMem  *phyMemory
	inst    inst
	cores   *multiCore
	console *console
	clicks  *clicks
	serial  *serial
	screen  *screen
	ticker  *ticker
	rom     *rom

	devices []device

	// Sections that are loaded into the machine
	Sections []*image.Section
}

// Default SP settings.
const (
	DefaultSPBase   uint32 = 0x20000
	DefaultSPStride uint32 = 0x2000
)

// NewMachine creates a machine with memory and cores.
// 0 memSize for full 4GB memory.
func NewMachine(c *Config) *Machine {
	if c.Ncore == 0 {
		c.Ncore = 1
	}
	ret := new(Machine)
	ret.phyMem = newPhyMemory(c.MemSize)
	ret.inst = new(instArch8)
	ret.cores = newMultiCore(c.Ncore, ret.phyMem, ret.inst)

	// hook-up devices
	p := ret.phyMem.Page(pageBasicIO)

	ret.console = newConsole(p, ret.cores)
	ret.clicks = newClicks(p, ret.cores)
	ret.serial = newSerial(p, ret.cores)
	ret.ticker = newTicker(ret.cores)

	ret.addDevice(ret.ticker)
	ret.addDevice(ret.serial)
	ret.addDevice(ret.console)
	ret.addDevice(ret.clicks)

	if c.Screen != nil {
		p1 := ret.phyMem.Page(pageScreenText)
		p2 := ret.phyMem.Page(pageScreenColor)
		ret.screen = newScreen(p1, p2, c.Screen)
		ret.addDevice(ret.screen)
	}

	sys := ret.phyMem.Page(pageSysInfo)
	sys.WriteWord(0, ret.phyMem.npage)
	sys.WriteWord(4, uint32(c.Ncore))

	if c.InitSP == 0 {
		ret.setSP(DefaultSPBase, DefaultSPStride)
	} else {
		ret.setSP(c.InitSP, c.StackPerCore)
	}
	ret.SetPC(c.InitPC)
	if c.Output != nil {
		ret.setOutput(c.Output)
	}
	if c.ROM != "" {
		ret.mountROM(c.ROM)
	}
	if c.RandSeed != 0 {
		ret.randSeed(c.RandSeed)
	}
	ret.setBootArg(c.BootArg)

	return ret
}

func (m *Machine) mountROM(root string) {
	p := m.phyMem.Page(pageBasicIO)
	m.rom = newROM(p, m.phyMem, m.cores, root)
	m.addDevice(m.rom)
}

func expToError(exp *Excep) error {
	if exp == nil {
		return nil
	}
	return exp
}

func (m *Machine) writePhyWord(phyAddr uint32, v uint32) error {
	return expToError(m.phyMem.WriteWord(phyAddr, v))
}

func (m *Machine) setBootArg(arg uint32) error {
	return m.writePhyWord(AddrBootArg, arg)
}

// ReadWord reads a word from the virtual address space.
func (m *Machine) ReadWord(core byte, virtAddr uint32) (uint32, error) {
	return m.cores.readWord(core, virtAddr)
}

// DumpRegs returns the values of the current registers of a core.
func (m *Machine) DumpRegs(core byte) []uint32 {
	return m.cores.dumpRegs(core)
}

func (m *Machine) setOutput(w io.Writer) {
	m.serial.Output = w
	m.console.Output = w
}

func (m *Machine) addDevice(d device) { m.devices = append(m.devices, d) }

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
			m.FlushScreen()
			return n, e
		}
	}

	return n, nil
}

// WriteBytes write a byte buffer to the memory at a particular offset.
func (m *Machine) WriteBytes(r io.Reader, offset uint32) error {
	start := offset % PageSize
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

func (m *Machine) randSeed(s int64) {
	m.ticker.Rand = rand.New(rand.NewSource(s))
}

func findCodeStart(secs []*image.Section) (uint32, bool) {
	for _, s := range secs {
		if s.Type == image.Code {
			return s.Addr, true
		}
	}
	return 0, false
}

// LoadSections loads a list of sections into the machine.
func (m *Machine) LoadSections(secs []*image.Section) error {
	for _, s := range secs {
		var buf io.Reader
		switch s.Type {
		case image.Zeros:
			buf = &zeroReader{s.Header.Size}
		case image.Code, image.Data:
			buf = bytes.NewReader(s.Bytes)
		case image.None, image.Debug, image.Comment:
			continue
		default:
			return fmt.Errorf("unknown section type: %d", s.Type)
		}

		if err := m.WriteBytes(buf, s.Addr); err != nil {
			return err
		}
	}

	if pc, found := findCodeStart(secs); found {
		m.SetPC(pc)
	}
	m.Sections = secs

	return nil
}

// SetPC sets all cores to start with a particular PC pointer.
func (m *Machine) SetPC(pc uint32) {
	for _, cpu := range m.cores.cores {
		cpu.regs[PC] = pc
	}
}

func (m *Machine) setSP(sp, stackSize uint32) {
	for i, cpu := range m.cores.cores {
		cpu.regs[SP] = sp + uint32(i+1)*stackSize
	}
}

// LoadImage loads an e8 image into the machine.
func (m *Machine) LoadImage(r io.ReadSeeker) error {
	secs, err := image.Read(r)
	if err != nil {
		return err
	}
	return m.LoadSections(secs)
}

// LoadImageBytes loads an e8 image in bytes into the machine.
func (m *Machine) LoadImageBytes(bs []byte) error {
	return m.LoadImage(bytes.NewReader(bs))
}

// PrintCoreStatus prints the cpu statuses.
func (m *Machine) PrintCoreStatus() { m.cores.PrintStatus() }

// FlushScreen flushes updates in the frame buffer to the
// screen device, even if the device has not asked for an update.
func (m *Machine) FlushScreen() {
	if m.screen != nil {
		m.screen.flush()
	}
}

// Click sends in a mouse click at the particular location.
func (m *Machine) Click(line, col uint8) { m.clicks.addClick(line, col) }

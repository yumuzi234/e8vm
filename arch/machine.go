package arch

import (
	"bytes"
	"fmt"
	"io"
	"math/rand"
	"time"

	"shanhu.io/smlvm/arch/devs"
	"shanhu.io/smlvm/image"
)

// Machine is a multicore shared memory simulated arch8 machine.
type Machine struct {
	phyMem *phyMemory
	inst   inst
	calls  *calls

	devices  []device
	console  *console
	clicks   *devs.ScreenClicks
	screen   *devs.Screen
	table    *devs.Table
	dialog   *devs.Dialog
	keyboard *devs.Keyboard
	rand     *devs.Rand
	ticker   *ticker
	rom      *rom

	cores *multiCore

	// Sections that are loaded into the machine
	sections []*image.Section
}

// Default SP settings.
const (
	DefaultSPBase   uint32 = 0x1000000
	DefaultSPStride uint32 = 0x2000
)

func makeRand(c *Config) *devs.Rand {
	if c.RandSeed == 0 {
		return devs.NewTimeRand()
	}
	return devs.NewRand(c.RandSeed)
}

// NewMachine creates a machine with memory and cores.
// 0 memSize for full 4GB memory.
func NewMachine(c *Config) *Machine {
	if c.Ncore == 0 {
		c.Ncore = 1
	}
	m := new(Machine)
	m.phyMem = newPhyMemory(c.MemSize)
	m.inst = new(instArch8)
	m.calls = newCalls(m.phyMem.Page(pageRPC), m.phyMem, c.Net)
	m.cores = newMultiCore(c.Ncore, m.phyMem, m.calls, m.inst)

	// hook-up devices
	p := m.phyMem.Page(pageBasicIO)

	m.console = newConsole(p, m.cores)
	m.ticker = newTicker(m.cores)

	m.calls.register(serviceConsole, m.console)
	m.calls.register(serviceRand, makeRand(c))
	now := time.Now()
	clk := &devs.Clock{
		Now:       c.Now,
		PerfNow:   c.PerfNow,
		StartTime: &now,
	}
	m.calls.register(serviceClock, clk)

	m.addDevice(m.ticker)
	m.addDevice(m.console)

	if c.Screen != nil {
		m.clicks = devs.NewScreenClicks(m.calls.sender(serviceScreen))
		s := devs.NewScreen(c.Screen)
		m.screen = s
		m.addDevice(s)
		m.calls.register(serviceScreen, s)
	}

	if c.Table != nil {
		t := devs.NewTable(c.Table, m.calls.sender(serviceTable))
		m.table = t
		m.calls.register(serviceTable, t) // hook vpc all
	}

	if c.Dialog != nil {
		d := devs.NewDialog(c.Dialog, m.calls.sender(serviceDialog))
		m.dialog = d
		m.calls.register(serviceDialog, d)
	}

	m.keyboard = devs.NewKeyboard(m.calls.sender(serviceKeyboard))

	sys := m.phyMem.Page(pageSysInfo)
	sys.WriteU32(0, m.phyMem.npage)
	sys.WriteU32(4, uint32(c.Ncore))

	if c.InitSP == 0 {
		m.cores.setSP(DefaultSPBase, DefaultSPStride)
	} else {
		m.cores.setSP(c.InitSP, c.StackPerCore)
	}
	m.cores.setPC(c.InitPC)
	if c.Output != nil {
		m.console.Output = c.Output
	}
	if c.ROM != "" {
		m.mountROM(c.ROM)
	}
	if c.RandSeed != 0 {
		m.randSeed(c.RandSeed)
	}
	m.phyMem.WriteU32(AddrBootArg, c.BootArg) // ignoring write error

	return m
}

func (m *Machine) mountROM(root string) {
	p := m.phyMem.Page(pageBasicIO)
	m.rom = newROM(p, m.phyMem, m.cores, root)
	m.addDevice(m.rom)
}

// ReadWord reads a word from the virtual address space.
func (m *Machine) ReadWord(core byte, virtAddr uint32) (uint32, error) {
	return m.cores.readWord(core, virtAddr)
}

// DumpRegs returns the values of the current registers of a core.
func (m *Machine) DumpRegs(core byte) []uint32 {
	return m.cores.dumpRegs(core)
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
	return m.phyMem.writeBytes(r, offset)
}

func (m *Machine) randSeed(s int64) {
	m.ticker.Rand = rand.New(rand.NewSource(s))
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

	if pc, found := image.CodeStart(secs); found {
		m.cores.setPC(pc)
	}
	m.sections = secs

	return nil
}

// LoadImage loads an smlvm image into the machine.
func (m *Machine) LoadImage(r io.ReadSeeker) error {
	secs, err := image.Read(r)
	if err != nil {
		return err
	}
	return m.LoadSections(secs)
}

// LoadImageBytes loads an smlvm image in bytes into the machine.
func (m *Machine) LoadImageBytes(bs []byte) error {
	return m.LoadImage(bytes.NewReader(bs))
}

// HandlePacket handles an incoming packet.
func (m *Machine) HandlePacket(p []byte) error {
	return m.calls.HandlePacket(p)
}

// PrintCoreStatus prints the cpu statuses.
func (m *Machine) PrintCoreStatus() { m.cores.PrintStatus() }

// FlushScreen flushes updates in the frame buffer to the
// screen device, even if the device has not asked for an update.
func (m *Machine) FlushScreen() {
	if m.screen != nil {
		m.screen.Flush()
	}
}

// Click sends in a mouse click at the particular location.
func (m *Machine) Click(line, col uint8) {
	if m.clicks == nil {
		return
	}
	m.clicks.Click(line, col)
}

// KeyDown sends in a key down event.
func (m *Machine) KeyDown(code uint8) { m.keyboard.KeyDown(code) }

// Choose sends in a user input choice.
func (m *Machine) Choose(index uint8) {
	if m.dialog == nil {
		return
	}
	m.dialog.Choose(index)
}

// ClickTable sends a click on the table at the particular location.
func (m *Machine) ClickTable(what string, pos uint8) {
	if m.table == nil {
		return
	}
	m.table.Click(what, pos)
}

// SleepTime returns the sleeping time required before next execution.
func (m *Machine) SleepTime() (time.Duration, bool) {
	return m.calls.sleepTime()
}

// HasPending checks if the machine has pending messages that are not
// delivered.
func (m *Machine) HasPending() bool {
	return m.calls.queueLen() > 0
}

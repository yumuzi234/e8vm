package arch

import (
	"fmt"
)

// MultiCore simulates a shared memory multicore processor.
type multiCore struct {
	cores  []*cpu
	phyMem *phyMemory
}

// NewMultiCore creates a shared memory multicore processor.
func newMultiCore(n int, mem *phyMemory, c *calls, i inst) *multiCore {
	if n > 32 {
		panic("too many cores")
	}

	ret := new(multiCore)
	ret.cores = make([]*cpu, n)
	ret.phyMem = mem

	for ind := range ret.cores {
		var thisCalls *calls
		if ind == 0 {
			thisCalls = c
		}
		ret.cores[ind] = newCPU(mem, thisCalls, i, byte(ind))
	}

	return ret
}

// CoreExcep is an exception on a particular core.
type CoreExcep struct {
	Core int
	*Excep
}

// Tick performs one tick on each core.
func (c *multiCore) Tick() *CoreExcep {
	for i, core := range c.cores {
		e := core.Tick()
		if e != nil {
			return &CoreExcep{i, e}
		}
	}

	return nil
}

// Ncore returns the number of cores.
func (c *multiCore) Ncore() byte {
	return byte(len(c.cores))
}

// Interrupt issues an interrupt to a particular core.
func (c *multiCore) Interrupt(code byte, core byte) {
	if int(core) >= len(c.cores) {
		panic("out of cores")
	}

	c.cores[core].Interrupt(code)
}

// PrintStatus prints out the core status of all the cores.
func (c *multiCore) PrintStatus() {
	for i, core := range c.cores {
		if len(c.cores) > 1 {
			fmt.Printf("[core %d]\n", i)
		}
		printCPUStatus(core)
		fmt.Println()
	}
}

func (c *multiCore) readWord(core byte, virtAddr uint32) (uint32, error) {
	if int(core) >= len(c.cores) {
		panic("out of cores")
	}

	v, exp := c.cores[core].virtMem.ReadU32(virtAddr, 0)
	if exp != nil {
		return 0, exp
	}
	return v, nil
}

func (c *multiCore) dumpRegs(core byte) []uint32 {
	if int(core) >= len(c.cores) {
		panic("out of cores")
	}

	ret := make([]uint32, Nreg)
	copy(ret, c.cores[core].regs)
	return ret
}

func (c *multiCore) setSP(sp, stackSize uint32) {
	for i, cpu := range c.cores {
		cpu.regs[SP] = sp + uint32(i+1)*stackSize
	}
}

func (c *multiCore) setPC(pc uint32) {
	for _, cpu := range c.cores {
		cpu.regs[PC] = pc
	}
}

func printCPUStatus(c *cpu) {
	p := func(name string, reg int) {
		fmt.Printf(" %3s = 0x%08x %-11d\n",
			name, c.regs[reg], int32(c.regs[reg]),
		)
	}

	p("r0", R0)
	p("r1", R1)
	p("r2", R2)
	p("r3", R3)
	p("r4", R4)
	p("sp", SP)
	p("ret", RET)
	p("pc", PC)

	fmt.Printf("ring = %d\n", c.ring)
}

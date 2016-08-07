package arch8

// InstSys exectues a system instruction
type instSys struct{}

func sysInfo(cpu *cpu, cmd uint32) (uint32, uint32) {
	switch cmd {
	case NCYCLE:
		ncycle := cpu.ncycle
		return uint32(ncycle), uint32(ncycle >> 16)
	case CPUID:
		return uint32(cpu.index), 0
	}

	return 0, 0
}

// I executes the system instruction.
// Returns any exception encountered.
func (i *instSys) I(cpu *cpu, in uint32) *Excep {
	op := (in >> 24) & 0xff // (32:24]
	r1 := (in >> 21) & 0x7  // (24:21]
	r2 := (in >> 18) & 0x7  // (21:18]
	v1 := cpu.regs[r1]
	v2 := cpu.regs[r2]

	switch op {
	case HALT:
		return errHalt
	case SYSCALL:
		if !cpu.UserMode() {
			return errInvalidInst
		}
		return cpu.Syscall()
	case JRUSER:
		cpu.ring = 1
		cpu.regs[PC] = v1
		return nil
	case VTABLE:
		if cpu.UserMode() {
			return errInvalidInst
		}
		cpu.virtMem.SetTable(v1)
	case IRET:
		if cpu.UserMode() {
			return errInvalidInst
		}
		return cpu.Iret()
	case SYSINFO:
		v1, v2 = sysInfo(cpu, v1)
	case SLEEP:
		if !cpu.sleeping {
			cpu.sleeping = true
			return errSleep
		}
		return nil
	default:
		return errInvalidInst
	}

	cpu.regs[r1] = v1
	cpu.regs[r2] = v2

	return nil
}

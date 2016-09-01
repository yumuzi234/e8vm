package dasm

import (
	"fmt"

	"e8vm.io/e8vm/arch8"
)

var (
	opSysMap = map[uint32]string{
		arch8.HALT:    "halt",
		arch8.SYSCALL: "syscall",
		arch8.IRET:    "iret",
		arch8.SLEEP:   "sleep",
	}

	opSys1Map = map[uint32]string{
		arch8.JRUSER: "jruser",
		arch8.VTABLE: "vtable",
	}

	opSys2Map = map[uint32]string{
		arch8.SYSINFO: "sysinfo",
	}
)

func instSys(addr uint32, in uint32) *Line {
	op := (in >> 24) & 0xff
	r1 := regStr((in >> 21) & 0x7)
	r2 := regStr((in >> 18) & 0x7)

	var s string
	if opStr, found := opSysMap[op]; found {
		s = opStr
	} else if opStr, found := opSys1Map[op]; found {
		s = fmt.Sprintf("%s %s", opStr, r1)
	} else if opStr, found := opSys2Map[op]; found {
		s = fmt.Sprintf("%s %s %s", opStr, r1, r2)
	}

	ret := newLine(addr, in)
	ret.Str = s
	return ret
}

package dasm

import (
	"fmt"

	"e8vm.io/e8vm/arch"
)

var (
	opShiftMap = map[uint32]string{
		arch.SLL: "sll",
		arch.SRL: "srl",
		arch.SRA: "sra",
	}

	opReg3Map = map[uint32]string{
		arch.SLLV: "sllv",
		arch.SRLV: "srlv",
		arch.SRLA: "srla",
		arch.ADD:  "add",
		arch.SUB:  "sub",
		arch.AND:  "and",
		arch.OR:   "or",
		arch.XOR:  "xor",
		arch.NOR:  "nor",
		arch.SLT:  "slt",
		arch.SLTU: "sltu",
		arch.MUL:  "mul",
		arch.MULU: "mulu",
		arch.DIV:  "div",
		arch.DIVU: "divu",
		arch.MOD:  "mod",
		arch.MODU: "modu",
	}

	opFloatMap = map[uint32]string{
		arch.FADD: "fadd",
		arch.FSUB: "fsub",
		arch.FMUL: "fmul",
		arch.FDIV: "fdiv",
		arch.FINT: "fint",
	}
)

func instReg(addr uint32, in uint32) *Line {
	if ((in >> 24) & 0xff) != 0 {
		panic("not a register inst")
	}

	dest := regStr((in >> 21) & 0x7)
	src1 := regStr((in >> 18) & 0x7)
	src2 := regStr((in >> 15) & 0x7)
	shift := (in >> 10) & 0x1f
	isFloat := (in >> 8) & 0x1
	funct := in & 0xff

	var s string
	if isFloat == 0 {
		if funct == arch.PANIC {
			s = fmt.Sprintf("panic")
		} else if funct == arch.SLLV && shift == 0 {
			s = fmt.Sprintf("mov %s %s", dest, src1)
		} else if opStr, found := opShiftMap[funct]; found {
			s = fmt.Sprintf("%s %s %s %d", opStr, dest, src1, shift)
		} else if opStr, found := opReg3Map[funct]; found {
			s = fmt.Sprintf("%s %s %s %s", opStr, dest, src1, src2)
		}
	} else {
		if opStr, found := opFloatMap[funct]; found {
			s = fmt.Sprintf("%s %s %s %s", opStr, dest, src1, src2)
		}
	}

	ret := newLine(addr, in)
	ret.Str = s

	return ret
}

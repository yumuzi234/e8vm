package codegen

import (
	A "e8vm.io/e8vm/arch8"
	S "e8vm.io/e8vm/asm8/inst"
)

// go lint (stupidly) forbids import with .
// so we will just copy the consts in here
const (
	_r0  = A.R0
	_r1  = A.R1
	_r2  = A.R2
	_r3  = A.R3
	_r4  = A.R4
	_ret = A.RET
	_sp  = A.SP
	_pc  = A.PC
)

// an empty struct for a separate namespace
type _s struct{}

var asm _s

func (_s) ims(op, d, s uint32, im int32) uint32 {
	return S.Imm(op, d, s, uint32(im))
}
func (_s) lw(d, s uint32, im int32) uint32 {
	return asm.ims(A.LW, d, s, im)
}
func (_s) sw(d, s uint32, im int32) uint32 {
	return asm.ims(A.SW, d, s, im)
}
func (_s) sb(d, s uint32, im int32) uint32 {
	return asm.ims(A.SB, d, s, im)
}
func (_s) lb(d, s uint32, im int32) uint32 {
	return asm.ims(A.LB, d, s, im)
}
func (_s) lbu(d, s uint32, im int32) uint32 {
	return asm.ims(A.LBU, d, s, im)
}
func (_s) addi(d, s uint32, im int32) uint32 {
	return asm.ims(A.ADDI, d, s, im)
}

func (_s) lui(d, im uint32) uint32     { return S.Imm(A.LUI, d, 0, im) }
func (_s) ori(d, s, im uint32) uint32  { return S.Imm(A.ORI, d, s, im) }
func (_s) xori(d, s, im uint32) uint32 { return S.Imm(A.XORI, d, s, im) }
func (_s) andi(d, s, im uint32) uint32 { return S.Imm(A.ANDI, d, s, im) }

func (_s) reg(op, d, s1, s2 uint32) uint32 {
	return S.Reg(op, d, s1, s2, 0, 0)
}

func (_s) add(d, s1, s2 uint32) uint32  { return asm.reg(A.ADD, d, s1, s2) }
func (_s) sub(d, s1, s2 uint32) uint32  { return asm.reg(A.SUB, d, s1, s2) }
func (_s) mul(d, s1, s2 uint32) uint32  { return asm.reg(A.MUL, d, s1, s2) }
func (_s) mulu(d, s1, s2 uint32) uint32 { return asm.reg(A.MULU, d, s1, s2) }
func (_s) div(d, s1, s2 uint32) uint32  { return asm.reg(A.DIV, d, s1, s2) }
func (_s) divu(d, s1, s2 uint32) uint32 { return asm.reg(A.DIVU, d, s1, s2) }
func (_s) mod(d, s1, s2 uint32) uint32  { return asm.reg(A.MOD, d, s1, s2) }
func (_s) modu(d, s1, s2 uint32) uint32 { return asm.reg(A.MODU, d, s1, s2) }
func (_s) and(d, s1, s2 uint32) uint32  { return asm.reg(A.AND, d, s1, s2) }
func (_s) or(d, s1, s2 uint32) uint32   { return asm.reg(A.OR, d, s1, s2) }
func (_s) xor(d, s1, s2 uint32) uint32  { return asm.reg(A.XOR, d, s1, s2) }
func (_s) nor(d, s1, s2 uint32) uint32  { return asm.reg(A.NOR, d, s1, s2) }
func (_s) sltu(d, s1, s2 uint32) uint32 { return asm.reg(A.SLTU, d, s1, s2) }
func (_s) slt(d, s1, s2 uint32) uint32  { return asm.reg(A.SLT, d, s1, s2) }
func (_s) sllv(d, s1, s2 uint32) uint32 { return asm.reg(A.SLLV, d, s1, s2) }
func (_s) srla(d, s1, s2 uint32) uint32 { return asm.reg(A.SRLA, d, s1, s2) }
func (_s) srlv(d, s1, s2 uint32) uint32 { return asm.reg(A.SRLV, d, s1, s2) }

func (_s) srl(d, s1, v uint32) uint32 {
	return S.Reg(A.SRL, d, s1, 0, v, 0)
}

func (_s) beq(s1, s2 uint32, im int32) uint32 {
	return S.Br(A.BEQ, s1, s2, im)
}
func (_s) bne(s1, s2 uint32, im int32) uint32 {
	return S.Br(A.BNE, s1, s2, im)
}

func (_s) jal(im int32) uint32 { return S.Jmp(A.JAL, im) }
func (_s) j(im int32) uint32   { return S.Jmp(A.J, im) }

func (_s) halt() uint32 { return S.Sys(A.HALT, 0) }

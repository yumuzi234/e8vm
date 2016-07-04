package arch8

// Immediate instructions
const (
	ADDI = 1
	SLTI = 2
	ANDI = 3
	ORI  = 4
	XORI = 5
	LUI  = 6

	LW  = 7
	LB  = 8
	LBU = 9
	SW  = 10
	SB  = 11
)

// Register instructions
const (
	PANIC = 0
	SLL   = 1
	SRL   = 2
	SRA   = 3
	SLLV  = 4
	SRLV  = 5
	SRLA  = 6
	ADD   = 7
	SUB   = 8
	AND   = 9
	OR    = 10
	XOR   = 11
	NOR   = 12
	SLT   = 13
	SLTU  = 14
	MUL   = 15
	MULU  = 16
	DIV   = 17
	DIVU  = 18
	MOD   = 19
	MODU  = 20

	FADD = 0
	FSUB = 1
	FMUL = 2
	FDIV = 3
	FINT = 4
)

// Branch instructions
const (
	BNE = 32
	BEQ = 33
)

// System instructions
const (
	HALT    = 64
	SYSCALL = 65
	JRUSER  = 66
	VTABLE  = 67
	IRET    = 68
	SYSINFO = 69
	SLEEP   = 70
)

// Jump instructions
const (
	J   = 2
	JAL = 3
)

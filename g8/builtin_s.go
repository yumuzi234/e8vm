package g8

const builtInSrc = `
// a char is sent in via r1
func PrintChar {
	// use r2 and r3
	addi sp sp -8
	sw r2 sp
	sw r3 sp 4

	ori r2 r0 0x2000 // the address of serial port
.wait
	lbu r3 r2 1
	bne r3 r0 .wait // wait for invalid

	mov r3 r1
	ori r3 r3 0x100
	sw r3 r2

	// restore r2 and r3
	lw r2 sp
	lw r3 sp 4

	addi sp sp 8
	mov pc ret
}

// Print a 32-bit unsigned integer
// when array is implemented, this will be rewritten in glang
func PrintUint32 {
	// saving used registers
	sw ret sp -4
	addi sp sp -28
	sw r1 sp
	sw r2 sp 4
	sw r3 sp 8
	
	bne r1 r0 .nonzero
.zero
	addi r1 r0 0x30 // '0'
	jal PrintChar
	j .end

.nonzero
    addi r2 sp 12
    ori r4 r0 10

.divloop
    modu r3 r1 r4
    sb r3 r2 0
    divu r1 r1 r4
    beq r1 r0 .print
    addi r2 r2 1
    j .divloop

.print
    addi r3 sp 12 // base

.printloop
    lbu r1 r2 0 // load
    addi r1 r1 0x30
    jal PrintChar
    beq r3 r2 .end
    addi r2 r2 -1
    j .printloop

.end
	addi r1 r0 0xa
	jal PrintChar // print a end line

	lw r2 sp 4
	lw r3 sp 8
	addi sp sp 28
	lw pc sp -4
}

// Print a 32-bit signed integer
// when array is implemented, this will be rewritten in glang
func PrintInt32 {
	// saving used registers
	sw ret sp -4
	addi sp sp -16 
	sw r1 sp
	sw r2 sp 4
	sw r3 sp 8
	
    slt r2 r1 r0 // r2 = r1 < 0
    beq r2 r0 .skipsign

    addi r1 r0 0x2d // '-'
    jal PrintChar

    lw r1 sp
    sub r1 r0 r1 // revert
.skipsign
    jal PrintUint32

	lw r2 sp 4
	lw r3 sp 8
	addi sp sp 16
	lw pc sp -4
}

// Panic halts the system immediately with panic exception
func Panic {
	panic
}

// Halt halts the system with halt exception
func Halt {
    halt
}

// MemCopy copies a range of memory
//    r1 - destination address
//    r2 - source address
//    r3 - number of bytes
func MemCopy {
	beq r3 r0 .ret
	beq r1 r2 .ret

	sltu r4 r1 r2
	beq r4 r0 .rev

.loop
	lbu r4 r2
	sb r4 r1
	addi r1 r1 1
	addi r2 r2 1
	addi r3 r3 -1
	bne r3 r0 .loop
	j .ret

.rev
	add r1 r1 r3
	add r2 r2 r3

.revloop
	addi r1 r1 -1
	addi r2 r2 -1
	lbu r4 r2
	sb r4 r1
	addi r3 r3 -1
	bne r3 r0 .revloop

.ret
	mov pc ret
}

// MemSet sets a range of memory to the same byte
//    r1 - destination address
//    r2 - the byte
//    r3 - number of bytes
func MemSet {
	beq r3 r0 .ret
.loop
	sb r2 r1
	addi r1 r1 1
	addi r3 r3 -1
	bne r3 r0 .loop
.ret
	mov pc ret
}

// MemClear sets a range of memory to zero
//   r1 - destination address
//   r2 - number of bytes
func MemClear {
	beq r2 r0 .ret
.loop
	sb r0 r1
	addi r1 r1 1
	addi r2 r2 -1
	bne r2 r0 .loop
.ret
	mov pc ret
}

func Vtable {
	vtable r1
	mov pc ret
}

func Syscall {
	syscall
	mov pc ret
}
`

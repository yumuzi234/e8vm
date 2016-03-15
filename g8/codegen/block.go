package codegen

import (
	"fmt"
)

const (
	jmpAlways = iota
	jmpIf
	jmpIfNot
)

type blockJump struct {
	typ  int
	cond Ref
	to   *Block
}

// Block is a basic block
type Block struct {
	id  int
	ops []Op

	insts    []*inst
	jumpInst *inst
	spMoved  bool

	frameSize *int32

	jump *blockJump

	next *Block // next in the linked list

	instStart int32
	instEnd   int32
}

func checkRefs(refs ...Ref) {
	for _, ref := range refs {
		if ref == nil {
			panic("nil ref")
		}
	}
}

func (b *Block) String() string { return fmt.Sprintf("B%d", b.id) }

func (b *Block) addOp(op Op) { b.ops = append(b.ops, op) }

// Comment adds an IR comment.
func (b *Block) Comment(s string) {
	b.addOp(&Comment{s})
}

// Commentf adds an IR comment with a particular printing format.
func (b *Block) Commentf(s string, args ...interface{}) {
	b.Comment(fmt.Sprintf(s, args...))
}

// Arith append an arithmetic operation to the basic block
func (b *Block) Arith(dest Ref, x Ref, op string, y Ref) {
	b.addOp(&ArithOp{dest, x, op, y})
}

// Assign appends an assignment operation to the basic block
func (b *Block) Assign(dest Ref, src Ref) {
	b.Arith(dest, nil, "", src)
}

// Zero appends zeroing operation to the basic block
func (b *Block) Zero(dest Ref) {
	checkRefs(dest)
	b.Arith(dest, nil, "0", nil)
}

// Call appends a function call operation to the basic block
func (b *Block) Call(dests []Ref, f Ref, args ...Ref) {
	checkRefs(dests...)
	checkRefs(f)
	checkRefs(args...)

	argsCopy := make([]Ref, len(args))
	copy(argsCopy, args)
	b.addOp(&CallOp{dests, f, argsCopy})
}

// Jump sets the block always jump to the dest block at its end
func (b *Block) Jump(dest *Block) {
	if dest == b.next {
		b.jump = nil
	} else {
		b.jump = &blockJump{jmpAlways, nil, dest}
	}
}

// JumpIfNot sets the block to jump to its natural next when the
// condition is met, and jump to dest when the condition is not met
func (b *Block) JumpIfNot(cond Ref, dest *Block) {
	b.jump = &blockJump{jmpIfNot, cond, dest}
}

// JumpIf sets the block to jump to its natural next when the
// condition is not met, and jump to dest when the condition is met
func (b *Block) JumpIf(cond Ref, dest *Block) {
	b.jump = &blockJump{jmpIf, cond, dest}
}

func (b *Block) inst(i uint32) *inst {
	ret := &inst{inst: i}
	b.insts = append(b.insts, ret)
	return ret
}

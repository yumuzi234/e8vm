package ir

// JumpType is the jump type
type JumpType int

// Jump types
const (
	JumpAlways JumpType = iota
	JumpIf
	JumpIfNot
)

// Jump is the jumping structure of
type Jump struct {
	Type JumpType
	Cond Ref
	To   *Block
}

// Block is a basic block
type Block struct {
	ID   int
	Ops  []Op
	Jump *Jump
	Next *Block
}

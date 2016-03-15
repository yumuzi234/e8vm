package ir

// Op is a general IR operation
type Op interface{}

// TacOp is a three-address-code operation
type TacOp struct {
	Dest Ref
	A    Ref
	Op   string
	B    Ref
}

// CallOp is a function call operation
type CallOp struct {
	Dest []Ref
	F    Ref
	Args []Ref
}

// Comment is a comment line for debugging
type Comment struct {
	Str string
}

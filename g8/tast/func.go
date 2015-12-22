package tast

// Func is a function.
type Func struct {
	Recv *Ref // the function receiver

	Body *Block
}

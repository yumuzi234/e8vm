package types

import (
	"bytes"
	"fmt"

	"e8vm.io/e8vm/arch8"
	"e8vm.io/e8vm/fmt8"
	"e8vm.io/e8vm/g8/ir"
)

// Arg is a function argument or return value
type Arg struct {
	Name string // optional
	T
}

// String returns "name T" for the named argument and "T" for an
// anonymous argument
func (a *Arg) String() string {
	if a.Name == "" {
		return a.T.String()
	}
	return fmt.Sprintf("%s %s", a.Name, a.T)
}

// Func is a function pointer type.
// It represents a particular function signature in G language.
type Func struct {
	Args     []*Arg
	Rets     []*Arg
	RetTypes []T
	Sig      *ir.FuncSig // the signature for IR

	// The method function signature.
	MethodFunc *Func

	// If the function pointer has a this pointer bond to it.
	IsBond bool
}

// NewFunc creates a new function type
func NewFunc(this *Arg, args []*Arg, rets []*Arg) *Func {
	ret := new(Func)

	if this != nil {
		ret.Args = append(ret.Args, this)
	}
	ret.Args = append(ret.Args, args...)
	ret.Rets = rets

	ret.Sig = makeFuncSig(ret)
	ret.RetTypes = argTypes(ret.Rets)

	if this != nil {
		ret.MethodFunc = NewFunc(nil, args, rets)
		ret.MethodFunc.IsBond = true
	}
	return ret
}

func argTypes(args []*Arg) []T {
	if args == nil {
		return nil
	}
	ret := make([]T, 0, len(args))
	for _, arg := range args {
		ret = append(ret, arg.T)
	}
	return ret
}

// NewFuncUnamed creates a new function type where all its arguments
// and return values are anonymous.
func NewFuncUnamed(args []T, rets []T) *Func {
	f := new(Func)
	for _, arg := range args {
		f.Args = append(f.Args, &Arg{T: arg})
	}
	for _, ret := range rets {
		f.Rets = append(f.Rets, &Arg{T: ret})
	}

	f.Sig = makeFuncSig(f)
	f.RetTypes = rets
	return f
}

// NewVoidFunc creates a new function that does not return anything.
func NewVoidFunc(args ...T) *Func { return NewFuncUnamed(args, nil) }

// VoidFunc is the signature for "func main()"
var VoidFunc = NewVoidFunc()

// String returns the function signature (without the argument names).
func (t *Func) String() string {
	// TODO: this is kind of ugly, need some refactor
	buf := new(bytes.Buffer)
	fmt.Fprintf(buf, "func (%s)", fmt8.Join(t.Args, ","))
	if len(t.Rets) > 1 {
		fmt.Fprintf(buf, " (")
		for i, ret := range t.Rets {
			if i > 0 {
				fmt.Fprint(buf, ",")
			}
			fmt.Fprint(buf, ret)
		}
		fmt.Fprint(buf, ")")
	} else if len(t.Rets) == 1 {
		fmt.Fprint(buf, " ")
		fmt.Fprint(buf, t.Rets[0])
	}

	if t.IsBond {
		fmt.Fprint(buf, " (bond)")
	}

	return buf.String()
}

// Size returns the size of a function pointer,
// which is equivalent to the size of a register.
func (t *Func) Size() int32 { return arch8.RegSize }

// RegSizeAlign returns true. Function pointer is always word aligned.
func (t *Func) RegSizeAlign() bool { return true }

func makeArg(t *Arg) *ir.FuncArg {
	return &ir.FuncArg{
		Name:         t.Name,
		Size:         t.Size(),
		U8:           IsBasic(t.T, Uint8),
		RegSizeAlign: t.RegSizeAlign(),
	}
}

// converts a langauge function signature into a IR function signature
func makeFuncSig(f *Func) *ir.FuncSig {
	narg := len(f.Args)
	args := make([]*ir.FuncArg, 0, narg)

	for _, t := range f.Args {
		if t.T == nil {
			panic("type missing")
		}
		args = append(args, makeArg(t))
	}

	rets := make([]*ir.FuncArg, len(f.Rets))
	for i, t := range f.Rets {
		rets[i] = makeArg(t)
	}

	return ir.NewFuncSig(args, rets)
}

// IsFuncPointer checks if the function is a simple function pointer
// that does not have a bond this pointer.
func IsFuncPointer(t T) bool {
	ft, ok := t.(*Func)
	if !ok {
		return false
	}

	return !ft.IsBond
}

package arch8

import (
	"errors"
	"fmt"
)

// Excep defines an exception error with a code
type Excep struct {
	Code byte
	Arg  uint32
	Err  error
}

// NewExcep creates a new Exception with a particular code and message.
func newExcep(c byte, s string) *Excep {
	ret := new(Excep)
	ret.Code = c
	ret.Err = errors.New(s)
	return ret
}

func (e *Excep) Error() string {
	if e.Arg != 0 {
		return fmt.Sprintf("%s: arg=%08x", e.Err.Error(), e.Arg)
	}
	return e.Err.Error()
}

// Exception codes
const (
	ErrHalt         = 1
	ErrTimer        = 2
	ErrInvalidInst  = 3
	ErrOutOfRange   = 4
	ErrMisalign     = 5
	ErrPageFault    = 6
	ErrPageReadonly = 7
	ErrPanic        = 8

	IntSerial = 16
	IntROM    = 17
	IntSwap   = 18
)

var (
	errHalt        = newExcep(ErrHalt, "halt")
	errTimeInt     = newExcep(ErrTimer, "time interrupt")
	errInvalidInst = newExcep(ErrInvalidInst, "invalid instruction")

	errMisalign = newExcep(ErrMisalign, "address misalign")
	errPanic    = newExcep(ErrPanic, "panic")
)

func newPageFault(va uint32) *Excep {
	ret := newExcep(ErrPageFault, "page fault")
	ret.Arg = va
	return ret
}

func newPageReadonly(va uint32) *Excep {
	ret := newExcep(ErrPageReadonly, "page read-only")
	ret.Arg = va
	return ret
}

func newOutOfRange(pa uint32) *Excep {
	ret := newExcep(ErrOutOfRange, "out of range")
	ret.Arg = pa
	return ret
}

package arch8

import (
	"bytes"
	"os"
)

// RunImageFile loads and run a raw image on a single core machine
// with 1GB physical memory until it runs into an exception.
func RunImageFile(path string) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}

	m := NewMachine(0, 1)
	if err := m.LoadImage(f); err != nil {
		return err
	}
	if _, exp := m.Run(0); exp != nil {
		return exp
	}

	return nil
}

func runImageArg(bs []byte, arg uint32, n int) (int, error) {
	m := NewMachine(0, 1)
	if err := m.LoadImageBytes(bs); err != nil {
		return 0, err
	}
	if err := m.WriteWord(AddrBootArg, arg); err != nil {
		return 0, err
	}

	ret, exp := m.Run(n)
	if exp == nil {
		return ret, nil
	}
	return ret, exp
}

// RunImage runs a series of bytes as a VM image with 1GB physical memory for
// maximum n cycles.  It returns the number of cycles, and the exit error if
// any.
func RunImage(bs []byte, n int) (int, error) {
	return runImageArg(bs, 0, n)
}

// RunImageArg runs a series of bytes as a VM image with 1GB physical memory
// until the machine shutsdown.  It returns the number of cycles, and the
// exit error if any.
func RunImageArg(bs []byte, arg uint32) (int, error) {
	return runImageArg(bs, arg, 0)
}

// RunImageOutput runs a image. It is similar to RunImage() but also returns
// the output.
func RunImageOutput(bs []byte, n int) (int, string, error) {
	m := NewMachine(0, 1)
	if err := m.LoadImageBytes(bs); err != nil {
		return 0, "", err
	}

	out := new(bytes.Buffer)
	m.SetOutput(out)

	ret, exp := m.Run(n)
	if exp == nil {
		return ret, out.String(), nil
	}
	return ret, out.String(), exp
}

// IsHalt returns true only when the error is a halt exception
func IsHalt(e error) bool { return IsErr(e, ErrHalt) }

// IsPanic returns true only when the error is a panic exception
func IsPanic(e error) bool { return IsErr(e, ErrPanic) }

// IsErr checks if the error is of a particular error code
func IsErr(e error, code byte) bool {
	if e, ok := e.(*Excep); ok {
		return e.Code == code
	}
	if e, ok := e.(*CoreExcep); ok {
		return e.Excep.Code == code
	}

	return false
}

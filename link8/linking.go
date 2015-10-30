package link8

import (
	"bytes"
	"fmt"
	"io"

	"e8vm.io/e8vm/arch8"
	"e8vm.io/e8vm/e8"
)

// Job is a linking job.
type Job struct {
	Pkg      *Pkg
	StartSym string
	InitPC   uint32
}

// NewJob creates a new linking job which init pc is the default one.
func NewJob(p *Pkg, start string) *Job {
	return &Job{
		Pkg:      p,
		StartSym: start,
		InitPC:   arch8.InitPC,
	}
}

// LinkMain is a short hand for NewJob(p, start).Link(out)
func LinkMain(p *Pkg, out io.Writer, start string) error {
	return NewJob(p, start).Link(out)
}

// Link performs the linking job and writes the output to out.
func (j *Job) Link(out io.Writer) error {
	lnk := newLinker()
	lnk.addPkgs(j.Pkg)

	var roots []uint32

	funcMain, index := j.Pkg.SymbolByName(j.StartSym)
	if funcMain == nil || funcMain.Type != SymFunc {
		return fmt.Errorf("start function missing")
	}
	roots = append(roots, index)
	used := traceUsed(lnk, j.Pkg, roots)

	funcs, vars, e := layout(used, j.InitPC)
	if e != nil {
		return e
	}

	buf := new(bytes.Buffer)

	w := newWriter(buf)
	for _, f := range funcs {
		writeFunc(w, f.pkg, f.Func())
	}
	for _, v := range vars {
		writeVar(w, v.pkg, v.Var())
	}
	if err := w.Err(); err != nil {
		return err
	}

	sec := &e8.Section{
		Header: &e8.Header{
			Type: e8.Code,
			Addr: j.InitPC,
		},
		Bytes: buf.Bytes(),
	}

	return e8.Write(out, []*e8.Section{sec})
}

// LinkBareFunc produces a image of a single function that has no links.
func LinkBareFunc(f *Func) ([]byte, error) {
	if f.TooLarge() {
		return nil, fmt.Errorf("code section too large")
	}

	buf := new(bytes.Buffer)
	w := newWriter(buf)
	w.writeBareFunc(f)
	if err := w.Err(); err != nil {
		return nil, err
	}

	image := new(bytes.Buffer)
	sec := &e8.Section{
		Header: &e8.Header{
			Type: e8.Code,
			Addr: arch8.InitPC,
		},
		Bytes: buf.Bytes(),
	}
	if err := e8.Write(image, []*e8.Section{sec}); err != nil {
		return nil, err
	}

	return image.Bytes(), nil
}

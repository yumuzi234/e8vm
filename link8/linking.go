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
	Pkgs     map[string]*Pkg
	Path     string
	StartSym string
	InitPC   uint32
}

// NewJob creates a new linking job which init pc is the default one.
func NewJob(pkgs map[string]*Pkg, path, start string) *Job {
	return &Job{
		Pkgs:     pkgs,
		Path:     path,
		StartSym: start,
		InitPC:   arch8.InitPC,
	}
}

// LinkMain is a short hand for NewJob(pkgs, path, start).Link(out)
func LinkMain(out io.Writer, pkgs map[string]*Pkg, path, start string) error {
	return NewJob(pkgs, path, start).Link(out)
}

// LinkSingle call LinkMain with only one single package.
func LinkSingle(out io.Writer, pkg *Pkg, start string) error {
	path := pkg.Path()
	pkgs := map[string]*Pkg{path: pkg}
	return LinkMain(out, pkgs, path, start)
}

// Link performs the linking job and writes the output to out.
func (j *Job) Link(out io.Writer) error {
	pkgs := j.Pkgs

	pkg := j.Pkgs[j.Path]

	funcMain := pkg.SymbolByName(j.StartSym)
	if funcMain == nil || funcMain.Type != SymFunc {
		return fmt.Errorf("start function missing")
	}

	roots := []string{j.StartSym}
	used := traceUsed(pkgs, j.Path, roots)

	funcs, vars, zeros, e := layout(used, j.InitPC)
	if e != nil {
		return e
	}

	var secs []*e8.Section
	if len(funcs) > 0 {
		buf := new(bytes.Buffer)
		w := newWriter(pkgs, buf)
		for _, f := range funcs {
			w.writeFunc(f.Func())
		}
		if err := w.Err(); err != nil {
			return err
		}

		if buf.Len() > 0 {
			secs = append(secs, &e8.Section{
				Header: &e8.Header{
					Type: e8.Code,
					Addr: j.InitPC,
				},
				Bytes: buf.Bytes(),
			})
		}
	}

	if len(vars) > 0 {
		buf := new(bytes.Buffer)
		w := newWriter(pkgs, buf)
		for _, v := range vars {
			w.writeVar(v.Var())
		}
		if err := w.Err(); err != nil {
			return err
		}

		if buf.Len() > 0 {
			secs = append(secs, &e8.Section{
				Header: &e8.Header{
					Type: e8.Data,
					Addr: vars[0].Var().addr,
				},
				Bytes: buf.Bytes(),
			})
		}
	}

	if len(zeros) > 0 {
		start := zeros[0].Var().addr
		lastVar := zeros[len(zeros)-1].Var()
		end := lastVar.addr + lastVar.Size()
		secs = append(secs, &e8.Section{
			Header: &e8.Header{
				Type: e8.Zeros,
				Addr: start,
				Size: end - start,
			},
		})
	}

	return e8.Write(out, secs)
}

// LinkBareFunc produces a image of a single function that has no links.
func LinkBareFunc(f *Func) ([]byte, error) {
	if f.TooLarge() {
		return nil, fmt.Errorf("code section too large")
	}

	buf := new(bytes.Buffer)
	w := newWriter(make(map[string]*Pkg), buf)
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

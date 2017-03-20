package builds

import (
	"path"
	"sort"
)

// Overlay is a input with two inputs.
type Overlay struct {
	in1 Input
	in2 Input
}

// NewOverlay creates a new input overlay.
func NewOverlay(in1, in2 Input) *Overlay {
	return &Overlay{
		in1: in1,
		in2: in2,
	}
}

// HasDir checks if the input has a directory.
func (o *Overlay) HasDir(p string) (bool, error) {
	ok, err := o.in1.HasDir(p)
	if err != nil {
		return false, err
	}
	if ok {
		return true, nil
	}

	return o.in2.HasDir(p)
}

// ListDirs lists all sub directories under a directory.
func (o *Overlay) ListDirs(p string) ([]string, error) {
	var ret []string
	retMap := make(map[string]bool)

	ok, err := o.in1.HasDir(p)
	if err != nil {
		return nil, err
	}

	if ok {
		dirs, err := o.in1.ListDirs(p)
		if err != nil {
			return nil, err
		}
		for _, dir := range dirs {
			retMap[dir] = true
			ret = append(ret, dir)
		}
	}

	ok, err = o.in1.HasDir(p)
	if err != nil {
		return nil, err
	}

	if ok {
		dirs, err := o.in2.ListDirs(p)
		if err != nil {
			return nil, err
		}

		for _, dir := range dirs {
			if retMap[dir] {
				continue
			}
			ret = append(ret, dir)
		}
	}

	sort.Strings(ret)
	return ret, nil
}

// ListFiles lists all files under a directory.
func (o *Overlay) ListFiles(p string) ([]string, error) {
	ok, err := o.in1.HasDir(p)
	if err != nil {
		return nil, err
	}
	if ok {
		return o.in1.ListFiles(p)
	}
	return o.in2.ListFiles(p)
}

// Open opens a file for reading.
func (o *Overlay) Open(p string) (*File, error) {
	dir := path.Dir(p)
	ok, err := o.in1.HasDir(dir)
	if err != nil {
		return nil, err
	}
	if ok {
		return o.in1.Open(p)
	}

	return o.in2.Open(p)
}

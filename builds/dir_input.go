package builds

import (
	"os"
	"path/filepath"
)

// DirInput is a input based on a directory.
type DirInput struct {
	dir string
}

// NewDirInput creates an input based on a file system directory.
func NewDirInput(dir string) *DirInput {
	if dir == "" {
		dir = "."
	}

	return &DirInput{dir: dir}
}

func (d *DirInput) p(p string) string {
	// TODO(h8liu): convert p to filepath)
	return filepath.Join(d.dir, p)
}

// HasDir chacks if the input has a directory.
func (d *DirInput) HasDir(p string) (bool, error) {
	info, err := os.Stat(d.p(p))
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, err
	}
	return info.IsDir(), nil
}

func (d *DirInput) readDir(p string) ([]os.FileInfo, error) {
	p = d.p(p)
	f, err := os.Open(p)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	infos, err := f.Readdir(0)
	if err != nil {
		return nil, err
	}

	if err := f.Close(); err != nil {
		return nil, err
	}
	return infos, nil
}

// ListDirs lists all sub directories under a directory.
func (d *DirInput) ListDirs(p string) ([]string, error) {
	infos, err := d.readDir(p)
	if err != nil {
		return nil, err
	}

	var ret []string
	for _, info := range infos {
		if info.IsDir() {
			ret = append(ret, info.Name())
		}
	}
	return ret, nil
}

// ListFiles lists all files under a directory.
func (d *DirInput) ListFiles(p string) ([]string, error) {
	infos, err := d.readDir(p)
	if err != nil {
		return nil, err
	}
	var ret []string
	for _, info := range infos {
		if !info.IsDir() {
			ret = append(ret, info.Name())
		}
	}
	return ret, nil
}

package builds

import (
	"io"
	"os"
	"path"
	"path/filepath"
)

// DirFS is a file system based on a directory.
type DirFS struct {
	dir string
}

// NewDirFS creates an input based on a file system directory.
func NewDirFS(dir string) *DirFS {
	if dir == "" {
		dir = "."
	}

	return &DirFS{dir: dir}
}

func (d *DirFS) p(p string) string {
	if p == "" {
		return d.dir
	}
	return filepath.Join(d.dir, filepath.FromSlash(p))
}

// HasDir chacks if the input has a directory.
func (d *DirFS) HasDir(p string) (bool, error) {
	info, err := os.Stat(d.p(p))
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, err
	}
	return info.IsDir(), nil
}

func (d *DirFS) readDir(p string) ([]os.FileInfo, error) {
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
func (d *DirFS) ListDirs(p string) ([]string, error) {
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
func (d *DirFS) ListFiles(p string) ([]string, error) {
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

// Open opens a file for reading.
func (d *DirFS) Open(p string) (*File, error) {
	name := path.Base(p)
	realPath := d.p(p)
	return &File{
		Name:   name,
		Path:   realPath,
		Opener: PathFile(realPath),
	}, nil
}

// Create creates a file for writing.
func (d *DirFS) Create(p string) (io.WriteCloser, error) {
	f, err := os.Open(d.p(p))
	if err != nil {
		return nil, err
	}
	return f, nil
}

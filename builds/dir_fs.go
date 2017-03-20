package builds

import (
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"
	"sort"
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

func (d *DirFS) hasDir(p string) (bool, error) {
	info, err := os.Stat(d.p(p))
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, err
	}
	return info.IsDir(), nil
}

// HasDir chacks if the input has a directory.
func (d *DirFS) HasDir(p string) (bool, error) {
	if err := checkValidDir(p); err != nil {
		return false, err
	}

	return d.hasDir(p)
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
	if err := checkValidDir(p); err != nil {
		return nil, err
	}

	ok, err := d.hasDir(p)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, fmt.Errorf("directory %q not exist", p)
	}

	infos, err := d.readDir(p)
	if err != nil {
		return nil, err
	}

	var ret []string
	for _, info := range infos {
		if !info.IsDir() {
			continue
		}
		name := info.Name()
		if !isValidPathName(name) {
			continue
		}
		ret = append(ret, name)
	}

	sort.Strings(ret)
	return ret, nil
}

// ListFiles lists all files under a directory.
func (d *DirFS) ListFiles(p string) ([]string, error) {
	if err := checkValidDir(p); err != nil {
		return nil, err
	}

	ok, err := d.hasDir(p)
	if err != nil {
		return nil, err
	}
	if ok {
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
		sort.Strings(ret)
		return ret, nil
	}

	return nil, fmt.Errorf("directory %q not exist", p)
}

// Open opens a file for reading.
func (d *DirFS) Open(p string) (*File, error) {
	if err := CheckValidPath(p); err != nil {
		return nil, err
	}

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
	if err := checkValidDir(p); err != nil {
		return nil, err
	}

	p = d.p(p)
	dir := filepath.Dir(p)
	if dir != "" {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return nil, err
		}
	}

	f, err := os.Create(p)
	if err != nil {
		return nil, err
	}
	return f, nil
}

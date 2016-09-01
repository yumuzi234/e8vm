// Package home8 provides the default build homing for building
// G language programs
package home8

import (
	"io"
	"io/ioutil"
	"path"
	"strings"

	"e8vm.io/e8vm/asm"
	"e8vm.io/e8vm/build8"
	"e8vm.io/e8vm/glang"
)

// Home provides the default building home.
type Home struct {
	home build8.Home

	path string
	std  string
}

// NewDirHome creates a new default home based on a particular directory.
func NewDirHome(path string, std string) *Home {
	lang := glang.Lang(false)
	dirHome := build8.NewDirHome(path, lang)
	dirHome.AddLang("asm", asm.Lang())

	return NewHome(std, dirHome)
}

// NewHome wraps an home with the specified std path.
func NewHome(std string, h build8.Home) *Home {
	if std == "" {
		std = "/smallrepo/std"
	}
	if !strings.HasPrefix(std, "/") {
		std = "/" + std
	}

	return &Home{
		std:  std,
		home: h,
	}
}

// AbsPath converts a possibly std path
func (h *Home) AbsPath(p string) string {
	if strings.HasPrefix(p, "/") {
		return p
	}
	return path.Join(h.std, p)
}

func (h *Home) dirPath(p string) string {
	abs := h.AbsPath(p)
	return strings.TrimPrefix(abs, "/")
}

// HasPkg checks if a package exists
func (h *Home) HasPkg(p string) bool {
	if p == "asm/builtin" {
		return true
	}
	return h.home.HasPkg(h.dirPath(p))
}

// Pkgs lists all the packages with a particular prefix
func (h *Home) Pkgs(prefix string) []string {
	prefix = h.dirPath(prefix)
	pkgs := h.home.Pkgs(prefix)
	var ret []string
	for _, p := range pkgs {
		p = strings.TrimPrefix("/"+p, h.std)
		ret = append(ret, p)
	}
	return ret
}

func builtinSrc() map[string]*build8.File {
	return map[string]*build8.File{
		"builtin.s": {
			Name:       "builtin.s",
			Path:       "<internal>/asm/builtin/builtin.s",
			ReadCloser: ioutil.NopCloser(strings.NewReader(glang.BuiltInSrc)),
		},
	}
}

// Src lists all the source files inside a package.
func (h *Home) Src(p string) map[string]*build8.File {
	if p == "asm/builtin" {
		return builtinSrc()
	}

	return h.home.Src(h.dirPath(p))
}

// Bin returns the wirter to write the binary image.
func (h *Home) Bin(p string) io.WriteCloser {
	return h.home.Bin(h.dirPath(p))
}

// TestBin returns the writer to write the test binary image.
func (h *Home) TestBin(p string) io.WriteCloser {
	return h.home.Bin(h.dirPath(p))
}

// Output returns the debug output writer for a particular name.
func (h *Home) Output(p, name string) io.WriteCloser {
	return h.home.Output(h.dirPath(p), name)
}

// Lang returns the langauge for the particular path.
// It returns assembly when any of the package name in the path
// is "asm".
func (h *Home) Lang(p string) build8.Lang {
	return h.home.Lang(h.dirPath(p))
}

// Package home8 provides the default build homing for building
// G language programs
package home8

import (
	"io"
	"path"
	"strings"

	"e8vm.io/e8vm/asm8"
	"e8vm.io/e8vm/build8"
	"e8vm.io/e8vm/g8"
)

// Home provides the default building home.
type Home struct {
	g       build8.Lang
	asm     build8.Lang
	dirHome *build8.DirHome

	path string
	std  string
}

// NewHome creates a new default home based on a particular directory.
func NewHome(path string, std string) *Home {
	if std == "" {
		std = "/smallrepo/std"
	}
	if strings.HasPrefix(std, "/") {
		std = "/" + std
	}

	lang := g8.Lang(false)
	dirHome := build8.NewDirHome(path, lang)
	dirHome.AddLang("asm", asm8.Lang())

	return &Home{
		g:       lang,
		asm:     asm8.Lang(),
		path:    path,
		std:     std,
		dirHome: dirHome,
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
	return h.dirHome.HasPkg(h.dirPath(p))
}

// Pkgs lists all the packages with a particular prefix
func (h *Home) Pkgs(prefix string) []string {
	prefix = h.dirPath(prefix)
	pkgs := h.dirHome.Pkgs(prefix)
	var ret []string
	for _, p := range pkgs {
		p = strings.TrimPrefix("/"+p, h.std)
		ret = append(ret, p)
	}
	return ret
}

type noopCloser struct{ io.Reader }

func (c *noopCloser) Close() error { return nil }

func builtinSrc() map[string]*build8.File {
	return map[string]*build8.File{
		"builtin.s": {
			Name: "builtin.s",
			Path: "<internal>/asm/builtin/builtin.s",
			ReadCloser: &noopCloser{
				Reader: strings.NewReader(g8.BuiltInSrc),
			},
		},
	}
}

// Src lists all the source files inside a package.
func (h *Home) Src(p string) map[string]*build8.File {
	if p == "asm/builtin" {
		return builtinSrc()
	}

	return h.dirHome.Src(h.dirPath(p))
}

// Bin returns the wirter to write the binary image.
func (h *Home) Bin(p string) io.WriteCloser {
	return h.dirHome.Bin(h.dirPath(p))
}

// TestBin returns the writer to write the test binary image.
func (h *Home) TestBin(p string) io.WriteCloser {
	return h.dirHome.Bin(h.dirPath(p))
}

// Output returns the debug output writer for a particular name.
func (h *Home) Output(p, name string) io.WriteCloser {
	return h.dirHome.Output(h.dirPath(p), name)
}

// Lang returns the langauge for the particular path.
// It returns assembly when any of the package name in the path
// is "asm".
func (h *Home) Lang(p string) build8.Lang {
	return h.dirHome.Lang(h.dirPath(p))
}

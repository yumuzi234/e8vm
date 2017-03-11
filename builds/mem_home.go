package builds

import (
	"io"
	"path"
	"sort"
)

// MemHome is a memory based building home.
type MemHome struct {
	pkgs  map[string]*MemPkg
	langs *LangPicker
}

// NewMemHome creates a new memory-based home
func NewMemHome(lang *Lang) *MemHome {
	return &MemHome{
		pkgs:  make(map[string]*MemPkg),
		langs: NewLangPicker(lang),
	}
}

// NewPkg creates (or replaces) a package of a particular path in this home.
func (h *MemHome) NewPkg(p string) *MemPkg {
	if !path.IsAbs(p) {
		p = path.Join("/", p)
	}

	ret := newMemPkg(p)
	h.pkgs[p] = ret
	return ret
}

// HasPkg checks if it has a particular package.
func (h *MemHome) HasPkg(p string) bool {
	_, found := h.pkgs[p]
	return found
}

// Pkgs lists all the packages in this home
func (h *MemHome) Pkgs(prefix string) []string {
	var ret []string
	for p := range h.pkgs {
		if IsParentPkg(prefix, p) {
			ret = append(ret, p)
		}
	}
	sort.Strings(ret)
	return ret
}

// Src lists all the source files in this home
func (h *MemHome) Src(p string) map[string]*File {
	pkg := h.pkgs[p]
	if pkg == nil {
		return nil
	}

	if len(pkg.files) == 0 {
		return nil
	}

	ret := make(map[string]*File)
	for name, f := range pkg.files {
		path := f.path
		if path == "" {
			path = "$" + p + "/" + name
		}
		ret[name] = &File{
			Path:   path,
			Name:   name,
			Opener: f.Opener(),
		}
	}

	return ret
}

// Bin opens the library binary for writing
func (h *MemHome) Bin(p string) io.WriteCloser {
	pkg := h.pkgs[p]
	if pkg == nil {
		panic("pkg not exists")
	}
	if pkg.bin == nil {
		pkg.bin = newMemFile("")
	} else {
		pkg.bin.Reset()
	}
	return pkg.bin
}

// TestBin opens the libary for writing the testing binary
func (h *MemHome) TestBin(p string) io.WriteCloser {
	pkg := h.pkgs[p]
	if pkg == nil {
		panic("pkg not exists")
	}
	if pkg.test == nil {
		pkg.test = newMemFile("")
	} else {
		pkg.test.Reset()
	}
	return pkg.test
}

// BinBytes returns the binary for the package if it has a main.
// It returns nil if the package does not.
// It panics if the package does not exist.
func (h *MemHome) BinBytes(p string) []byte {
	pkg := h.pkgs[p]
	if pkg == nil {
		panic("pkg not exists")
	}
	if pkg.bin == nil {
		return nil
	}
	return pkg.bin.Bytes()
}

// Output creates a debug output file for writing.
func (h *MemHome) Output(p, name string) io.WriteCloser {
	pkg := h.pkgs[p]
	if pkg == nil {
		panic("pkg not exists")
	}
	ret := newMemFile("")
	pkg.outs[name] = ret
	return ret
}

// OutputBytes gets the debug output bytes of a compilation.
func (h *MemHome) OutputBytes(p, name string) []byte {
	pkg := h.pkgs[p]
	if pkg == nil {
		panic("pkg not exists")
	}
	ret := pkg.outs[name]
	if ret == nil {
		return nil
	}
	return ret.Bytes()
}

// Lang returns the language for path
func (h *MemHome) Lang(path string) *Lang { return h.langs.Lang(path) }

// AddLang adds a language to a prefix
func (h *MemHome) AddLang(prefix string, lang *Lang) {
	h.langs.AddLang(prefix, lang)
}

// AddFiles adds a set of files into mem home.
func (h *MemHome) AddFiles(files map[string]string) {
	pkgs := make(map[string]*MemPkg)
	for f, content := range files {
		p := path.Dir(f)
		base := path.Base(f)
		pkg, found := pkgs[p]
		if !found {
			pkg = h.NewPkg(p)
		}
		pkg.AddFile(f, base, content)
	}
}

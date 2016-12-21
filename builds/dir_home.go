package builds

import (
	"io"
	"io/ioutil"
	"log"
	"os"
	"path"
	"path/filepath"
	"sort"
	"strings"

	"shanhu.io/smlvm/lexing"
)

func listSrcFiles(dir string, lang *Lang) ([]string, error) {
	files, e := ioutil.ReadDir(dir)
	if e != nil {
		return nil, e
	}

	var ret []string
	ext := "." + lang.Ext

	for _, file := range files {
		if file.IsDir() {
			continue
		}

		name := file.Name()
		if strings.HasSuffix(name, ext) {
			ret = append(ret, name)
		}
	}

	return ret, nil
}

// DirHome is a file system basd building home.
type DirHome struct {
	path  string
	langs *LangPicker

	fileList map[string][]string

	MemHome *MemHome
	Quiet   bool
}

// NewDirHome creates a file system home storage with
// a particualr default language for compiling.
func NewDirHome(path string, lang *Lang) *DirHome {
	if lang == nil {
		panic("must specify a default language")
	}

	if path == "" {
		path = "."
	}
	return &DirHome{
		path:     path,
		fileList: make(map[string][]string),
		langs:    NewLangPicker(lang),
	}
}

func (h *DirHome) out(pre, p string) string {
	return filepath.Join(h.path, "_", pre, p)
}

func (h *DirHome) outFile(pre, p, f string) string {
	return filepath.Join(h.path, "_", pre, p, f)
}

func (h *DirHome) srcFile(p, f string) string {
	return filepath.Join(h.path, p, f)
}

// ClearCache clears the file list cache
func (h *DirHome) ClearCache() {
	h.fileList = make(map[string][]string)
}

// AddLang registers a language with a particular path prefix
func (h *DirHome) AddLang(keyword string, lang *Lang) {
	h.langs.AddLang(keyword, lang)
}

// HasPkg checks if a package exists.
func (h *DirHome) HasPkg(p string) bool {
	if h.MemHome != nil {
		if h.MemHome.HasPkg(p) {
			return true
		}
	}

	root := h.path
	fp := filepath.Join(root, p)
	info, err := os.Stat(fp)
	if err != nil {
		return false
	}
	if !info.IsDir() {
		return false
	}

	base := filepath.Base(p)
	if !lexing.IsPkgName(base) {
		return false
	}

	lang := h.Lang(p)
	files, err := listSrcFiles(fp, lang)
	if err != nil {
		return false
	}

	if len(files) > 0 {
		h.fileList[p] = files // caching
		return true
	}

	return false
}

// Pkgs lists all the packages with a particular prefix.
func (h *DirHome) Pkgs(prefix string) []string {
	var pkgs []string
	if h.MemHome != nil {
		pkgs = h.MemHome.Pkgs(prefix)
	}

	root := h.path
	start := filepath.Join(root, prefix)

	walkFunc := func(p string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if p == "." || p == root {
			return nil
		}

		relPath, err := filepath.Rel(root, p)
		if err != nil {
			panic(err)
		}
		if relPath == "." {
			return nil
		}
		if !info.IsDir() {
			return nil
		}
		base := filepath.Base(relPath)
		if !lexing.IsPkgName(base) {
			return filepath.SkipDir
		}

		pkgPath := path.Join(filepath.SplitList(relPath)...)
		pkgPath = path.Join("/", pkgPath)

		lang := h.Lang(pkgPath)
		if lang == nil {
			panic(pkgPath)
		}

		files, err := listSrcFiles(p, lang)
		if err != nil {
			return err
		}

		if len(files) > 0 {
			h.fileList[pkgPath] = files // caching
			pkgs = append(pkgs, pkgPath)
		}

		return nil
	}

	e := filepath.Walk(start, walkFunc)
	if e != nil && !h.Quiet {
		log.Fatal("error", e)
	}

	sort.Strings(pkgs)
	return pkgs
}

// Src lists all the source files inside this package.
func (h *DirHome) Src(p string) map[string]*File {
	if !IsPkgPath(p) {
		panic("not package path")
	}

	if h.MemHome != nil && h.MemHome.HasPkg(p) {
		return h.MemHome.Src(p)
	}

	lang := h.Lang(p)
	if lang == nil {
		return nil
	}

	files, found := h.fileList[p]
	if !found {
		var err error
		files, err = listSrcFiles(p, lang)
		if err != nil {
			return nil
		}

		h.fileList[p] = files
	}

	if len(files) == 0 {
		return nil
	}

	ret := make(map[string]*File)
	for _, name := range files {
		filePath := h.srcFile(p, name)
		ret[name] = &File{
			Path:   filePath,
			Name:   name,
			Opener: PathFile(filePath),
		}
	}

	return ret
}

// Bin returns the writer to write the binary image.
func (h *DirHome) Bin(p string) io.WriteCloser {
	if !IsPkgPath(p) {
		panic("not package path")
	}
	if h.MemHome != nil && h.MemHome.HasPkg(p) {
		return h.MemHome.Bin(p)
	}
	return newDirFile(h.out("bin", p+".e8"))
}

// TestBin returns the writer to write the test binary image.
func (h *DirHome) TestBin(p string) io.WriteCloser {
	if !IsPkgPath(p) {
		panic("not package path")
	}
	if h.MemHome != nil && h.MemHome.HasPkg(p) {
		return h.MemHome.TestBin(p)
	}
	return newDirFile(h.out("test", p+".e8"))
}

// Output returns the debug output writer for a particular name.
func (h *DirHome) Output(p, name string) io.WriteCloser {
	if !IsPkgPath(p) {
		panic("not package path")
	}
	if h.MemHome != nil && h.MemHome.HasPkg(p) {
		return h.MemHome.Output(p, name)
	}
	return newDirFile(h.outFile("out", p, name))
}

// Lang returns the language for the particular path.
// It searches for the longest prefix match
func (h *DirHome) Lang(p string) *Lang {
	if h.MemHome != nil && h.MemHome.HasPkg(p) {
		return h.MemHome.Lang(p)
	}

	return h.langs.Lang(p)
}

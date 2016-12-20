package builds

import (
	"io"
	"io/ioutil"
	"log"
	"os"
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

	Quiet bool
}

// NewDirHome creates a file system home storage with
// a particualr default language for compiling.
func NewDirHome(path string, lang *Lang) *DirHome {
	if lang == nil {
		panic("must specify a default language")
	}

	ret := new(DirHome)
	ret.path = path
	ret.fileList = make(map[string][]string)
	ret.langs = NewLangPicker(lang)

	return ret
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

	files, err := listSrcFiles(p, lang)
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
	root := h.path
	start := filepath.Join(root, prefix)
	var pkgs []string

	walkFunc := func(p string, info os.FileInfo, e error) error {
		if e != nil {
			return e
		}
		if p == "." || p == root {
			return nil
		}

		path, e := filepath.Rel(root, p)
		if e != nil {
			panic(e)
		} else if path == "." {
			return nil
		}

		base := filepath.Base(path)
		if !info.IsDir() {
			return nil
		} else if !lexing.IsPkgName(base) {
			return filepath.SkipDir
		}

		lang := h.Lang(path)
		if lang == nil {
			panic(path)
		}

		files, err := listSrcFiles(p, lang)
		if err != nil {
			return err
		}

		if len(files) > 0 {
			h.fileList[path] = files // caching
			pkgs = append(pkgs, path)
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
	return newDirFile(h.out("bin", p+".e8"))
}

// TestBin returns the writer to write the test binary image.
func (h *DirHome) TestBin(p string) io.WriteCloser {
	if !IsPkgPath(p) {
		panic("not package path")
	}
	return newDirFile(h.out("test", p+".e8"))
}

// Output returns the debug output writer for a particular name.
func (h *DirHome) Output(p, name string) io.WriteCloser {
	if !IsPkgPath(p) {
		panic("not package path")
	}
	return newDirFile(h.outFile("out", p, name))
}

// Lang returns the language for the particular path.
// It searches for the longest prefix match
func (h *DirHome) Lang(p string) *Lang {
	return h.langs.Lang(p)
}

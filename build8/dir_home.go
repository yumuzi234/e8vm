package build8

import (
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"sort"

	"e8vm.io/e8vm/lex8"
)

func listSrcFiles(dir string, lang Lang) ([]string, error) {
	files, e := ioutil.ReadDir(dir)
	if e != nil {
		return nil, e
	}

	var ret []string

	for _, file := range files {
		if file.IsDir() {
			continue
		}

		name := file.Name()
		if lang.IsSrc(name) {
			ret = append(ret, name)
		}
	}

	return ret, nil
}

// DirHome is a file system basd building home.
type DirHome struct {
	path  string
	langs *langPicker

	fileList map[string][]string

	Quiet bool
}

// NewDirHome creates a file system home storage with
// a particualr default language for compiling.
func NewDirHome(path string, lang Lang) *DirHome {
	if lang == nil {
		panic("must specify a default language")
	}

	ret := new(DirHome)
	ret.path = path
	ret.fileList = make(map[string][]string)
	ret.langs = newLangPicker(lang)

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
func (h *DirHome) AddLang(prefix string, lang Lang) {
	h.langs.addLang(prefix, lang)
}

// Pkgs lists all the packages inside this home folder.
func (h *DirHome) Pkgs(prefix string) []string {
	root := h.path
	start := filepath.Join(root, prefix)
	var pkgs []string

	walkFunc := func(p string, info os.FileInfo, e error) error {
		if e != nil {
			return e
		}
		if p == "." {
			return nil
		}

		if !info.IsDir() {
			return nil
		} else if !lex8.IsPkgName(info.Name()) {
			return filepath.SkipDir
		}

		if root == p {
			return nil
		}

		path, e := filepath.Rel(root, p)
		if e != nil {
			panic(e)
		} else if path == "." {
			return nil
		}

		lang := h.Lang(path)
		if lang == nil {
			panic(path)
		}

		files, e := listSrcFiles(p, lang)
		if e != nil {
			return e
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
	if !isPkgPath(p) {
		panic("not package path")
	}

	lang := h.Lang(p)
	if lang == nil {
		return nil
	}

	files, found := h.fileList[p]
	if !found {
		files, e := listSrcFiles(p, lang)
		if e != nil {
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
			Path:       filePath,
			Name:       name,
			ReadCloser: newDirFile(filePath),
		}
	}

	return ret
}

// Bin returns the writer to write the binary image.
func (h *DirHome) Bin(p string) io.WriteCloser {
	if !isPkgPath(p) {
		panic("not package path")
	}
	return newDirFile(h.out("bin", p+".e8"))
}

// TestBin returns the writer to write the test binary image.
func (h *DirHome) TestBin(p string) io.WriteCloser {
	if !isPkgPath(p) {
		panic("not package path")
	}
	return newDirFile(h.out("test", p+".e8"))
}

// Output returns the debug output writer for the particular name.
func (h *DirHome) Output(p, name string) io.WriteCloser {
	if !isPkgPath(p) {
		panic("not package path")
	}
	return newDirFile(h.outFile("out", p, name))
}

// Lang returns the language for the particular path.
// It searches for the longest prefix match
func (h *DirHome) Lang(p string) Lang { return h.langs.lang(p) }

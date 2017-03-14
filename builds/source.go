package builds

import (
	"fmt"
	"path"
	"sort"
	"strings"
)

type source struct {
	in         Input2
	langPicker *LangPicker
}

func newSource(in Input2, langPicker *LangPicker) *source {
	return &source{
		in:         in,
		langPicker: langPicker,
	}
}

func langFiles(lang *Lang, files []string) []string {
	var ret []string
	ext := "." + lang.Ext
	for _, f := range files {
		if strings.HasSuffix(f, ext) {
			ret = append(ret, f)
		}
	}
	return ret
}

func (s *source) listSrcFiles(p string) ([]string, error) {
	files, err := s.in.ListFiles(p)
	if err != nil {
		return nil, err
	}
	lang := s.langPicker.Lang(p)
	return langFiles(lang, files), nil
}

func (s *source) srcFileMap(p string) (map[string]*File, error) {
	rel := relPath(p)
	ok, err := s.in.HasDir(rel)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, fmt.Errorf("directory %q not exist", rel)
	}

	files, err := s.listSrcFiles(rel)
	if err != nil {
		return nil, err
	}

	if len(files) == 0 {
		return nil, fmt.Errorf("package %q has no source file", p)
	}

	ret := make(map[string]*File)
	for _, name := range files {
		f, err := s.in.Open(path.Join(rel, name))
		if err != nil {
			return nil, err
		}
		ret[name] = f
	}

	return ret, nil
}

func (s *source) hasPkg(p string) (bool, error) {
	rel := relPath(p)
	ok, err := s.in.HasDir(rel)
	if err != nil {
		return false, err
	}
	if !ok {
		return false, nil
	}

	files, err := s.listSrcFiles(rel)
	if err != nil {
		return false, err
	}
	return len(files) > 0, nil
}

func (s *source) lang(p string) *Lang {
	return s.langPicker.Lang(p)
}

func (s *source) listPkgs(lst []string, p string) ([]string, error) {
	dirs, err := s.in.ListDirs(p)
	if err != nil {
		return nil, err
	}

	var subs []string
	for _, dir := range dirs {
		if strings.HasPrefix(dir, "_") {
			continue
		}
		subs = append(subs, path.Join(p, dir))
	}

	for _, sub := range subs {
		files, err := s.listSrcFiles(sub)
		if err != nil {
			return nil, err
		}
		if len(files) > 0 {
			lst = append(lst, path.Join("/", sub))
		}
		lst, err = s.listPkgs(lst, sub)
		if err != nil {
			return nil, err
		}
	}

	return lst, nil
}

func (s *source) allPkgs(p string) ([]string, error) {
	rel := relPath(p)
	ok, err := s.in.HasDir(rel)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, fmt.Errorf("directory %q not exist", rel)
	}
	var ret []string
	if rel != "" {
		ret = append(ret, rel)
	}
	ret, err = s.listPkgs(ret, rel)
	if err != nil {
		return nil, err
	}
	sort.Strings(ret)
	return ret, nil
}

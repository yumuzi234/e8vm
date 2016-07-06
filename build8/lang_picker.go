package build8

import (
	"strings"
)

type langPicker struct {
	defaultLang Lang
	langs       map[string]Lang
}

func newLangPicker(def Lang) *langPicker {
	if def == nil {
		panic("default language must not be nil")
	}

	ret := new(langPicker)
	ret.defaultLang = def
	ret.langs = make(map[string]Lang)
	return ret
}

func (pick *langPicker) addLang(key string, lang Lang) {
	if lang == nil {
		panic("language must not be nil")
	}

	if key == "" {
		pick.defaultLang = lang
	}
	pick.langs[key] = lang
}

func (pick *langPicker) lang(path string) Lang {
	if !isPkgPath(path) {
		panic("not package path")
	}

	pkgs := strings.Split(path, "/")
	for _, pkg := range pkgs {
		if ret, ok := pick.langs[pkg]; ok {
			return ret
		}
	}

	return pick.defaultLang
}

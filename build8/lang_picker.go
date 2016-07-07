package build8

import (
	"strings"
)

// LangPicker is a utility that picks a language based on a keyword.
type LangPicker struct {
	defaultLang Lang
	langs       map[string]Lang
}

// NewLangPicker creates a new picker with a default language.
func NewLangPicker(def Lang) *LangPicker {
	if def == nil {
		panic("default language must not be nil")
	}

	return &LangPicker{
		defaultLang: def,
		langs:       make(map[string]Lang),
	}
}

// AddLang adds a language that has a certain keyword.
// When key is empty, it replaces the default language.
func (pick *LangPicker) AddLang(key string, lang Lang) {
	if lang == nil {
		panic("language must not be nil")
	}

	if key == "" {
		pick.defaultLang = lang
	}
	pick.langs[key] = lang
}

// Lang picks the language for a particular path.
func (pick *LangPicker) Lang(path string) Lang {
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

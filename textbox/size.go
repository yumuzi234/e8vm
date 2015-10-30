package textbox

import (
	"bufio"
	"io"
	"os"
	"strings"

	"e8vm.io/e8vm/lex8"
	"e8vm.io/tools/goload"
)

// TabSize is the indent size for each tab
const TabSize = 4

func runeWidth(r rune) int {
	switch r {
	case '\t':
		return TabSize
	case '\n', '\r':
		return 0
	}
	return 1
}

// Rect returns the display of a text line.
// Ends of lines are ignored.
func Rect(r io.Reader) (nline, maxWidth int, e error) {
	br := bufio.NewReader(r)
	nline = 0
	curWidth := 0
	maxWidth = 0

	for {
		r, _, e := br.ReadRune()
		if e == io.EOF {
			if curWidth > 0 {
				nline++
			}
			break
		} else if e != nil {
			return 0, 0, e
		}

		if r == '\n' {
			nline++
			if curWidth > maxWidth {
				maxWidth = curWidth
			}
			curWidth = 0
		} else {
			curWidth += runeWidth(r)
		}
	}

	if curWidth > maxWidth {
		maxWidth = curWidth
	}

	return nline, maxWidth, nil
}

func checkRect(log lex8.Logger, file string, r io.Reader, h, w int) {
	br := bufio.NewReader(r)
	row := 0
	col := 0

	pos := func() *lex8.Pos { return &lex8.Pos{file, row + 1, col + 1} }
	newLine := func() {
		if col > w {
			log.Errorf(pos(), "line too wide")
		}
		row++
		col = 0
	}

	for {
		r, _, e := br.ReadRune()
		if e == io.EOF {
			if col > 0 {
				newLine()
			}
			break
		} else if lex8.LogError(log, e) {
			break
		}

		if r == '\n' {
			newLine()
		} else {
			col += runeWidth(r)
		}
	}

	if row > h && !strings.HasSuffix(file, "_bytes.go") {
		log.Errorf(pos(), "too many lines")
	}
}

// CheckRectLoaded checks all the golang source files in a loaded go program
func CheckRectLoaded(prog *goload.Program, h, w int) []*lex8.Error {
	errs := lex8.NewErrorList()

	fset := prog.Fset
	for _, p := range prog.Pkgs {
		pinfo := prog.Imported[p]
		for _, astf := range pinfo.Files {
			tokFile := fset.File(astf.Pos())
			name := tokFile.Name()
			fin, e := os.Open(name)
			if lex8.LogError(errs, e) {
				continue
			}

			checkRect(errs, name, fin, h, w)
			if lex8.LogError(errs, fin.Close()) {
				continue
			}
		}
	}

	return errs.Errs()
}

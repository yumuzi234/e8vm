package textbox

import (
	"bufio"
	"io"

	"shanhu.io/smlvm/lexing"
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

// CheckRect checks if a program is within a rectangular area.
func CheckRect(file string, r io.Reader, h, w int) []*lexing.Error {
	br := bufio.NewReader(r)
	row := 0
	col := 0

	errs := lexing.NewErrorList()

	pos := func() *lexing.Pos { return &lexing.Pos{file, row + 1, col + 1} }
	newLine := func() {
		if col > w {
			errs.Errorf(pos(), "line too wide")
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
		} else if lexing.LogError(errs, e) {
			break
		}

		if r == '\n' {
			newLine()
		} else {
			col += runeWidth(r)
		}
	}

	if row > h {
		errs.Errorf(pos(), "too many lines")
	}

	return errs.Errs()
}

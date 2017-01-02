package gfmt

import (
	"shanhu.io/smlvm/lexing"
	"shanhu.io/smlvm/pl/ast"
)

func printExprList(f *formatter, begin, end *lexing.Token, list *ast.ExprList) {
	pos := begin.Pos

	f.Tab()

	n := 0
	for i, e := range list.Exprs {
		p := ast.ExprPos(e)
		if p.Line > pos.Line {
			if p.Line > pos.Line+1 && i > 0 {
				f.printEndl()
			}
			f.printEndl()
			n = 0
		}

		if n > 0 {
			f.printExprs(" ")
		}

		f.printExprs(e)
		if i < len(list.Commas) {
			f.printExprs(list.Commas[i])
		}
		pos = p
		n++
	}

	if end.Pos.Line > pos.Line {
		f.printEndl()
	}

	f.ShiftTab()
}

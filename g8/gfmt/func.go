package gfmt

import (
	"fmt"

	"e8vm.io/e8vm/fmt8"
	"e8vm.io/e8vm/g8/ast"
)

func printParaList(p *fmt8.Printer, m *matcher, lst *ast.ParaList) {
	printToken(p, m, lst.Lparen)
	for i, para := range lst.Paras {
		if i > 0 {
			printExprs(p, m, lst.Commas[i-1], " ")
		}
		if para.Ident != nil {
			printToken(p, m, para.Ident)
			if para.Type != nil {
				fmt.Fprint(p, " ")
			}
		}

		if para.Type != nil {
			printExpr(p, m, para.Type)
		}
	}
	printToken(p, m, lst.Rparen)
}

func printFunc(p *fmt8.Printer, m *matcher, f *ast.Func) {
	printExprs(p, m, f.Kw, " ", f.Name)
	printParaList(p, m, f.Args)
	if f.RetType != nil {
		printExprs(p, m, " ", f.RetType)
	} else if f.Rets != nil {
		fmt.Fprint(p, " ")
		printParaList(p, m, f.Rets)
	}

	fmt.Fprint(p, " ")
	printStmt(p, m, f.Body)
	fmt.Fprintln(p)
}

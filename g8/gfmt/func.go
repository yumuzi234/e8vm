package gfmt

import (
	"e8vm.io/e8vm/g8/ast"
)

func printParaList(f *formatter, lst *ast.ParaList) {
	f.printToken(lst.Lparen)
	for i, para := range lst.Paras {
		if i > 0 {
			printExprs(f, lst.Commas[i-1], " ")
		}
		if para.Ident != nil {
			f.printToken(para.Ident)
			if para.Type != nil {
				f.printSpace()
			}
		}

		if para.Type != nil {
			printExpr(f, para.Type)
		}
	}
	f.printToken(lst.Rparen)
}

func printFunc(f *formatter, fn *ast.Func) {
	printExprs(f, fn.Kw, " ", fn.Name)
	printParaList(f, fn.Args)
	if fn.RetType != nil {
		printExprs(f, " ", fn.RetType)
	} else if fn.Rets != nil {
		f.printSpace()
		printParaList(f, fn.Rets)
	}

	f.printSpace()
	printStmt(f, fn.Body)
}

package gfmt

import (
	"e8vm.io/e8vm/g8/ast"
)

func printParaList(f *formatter, lst *ast.ParaList) {
	f.printToken(lst.Lparen)
	for i, para := range lst.Paras {
		if i > 0 {
			f.printExprs(lst.Commas[i-1], " ")
		}
		if para.Ident != nil {
			f.printToken(para.Ident)
			if para.Type != nil {
				f.printSpace()
			}
		}

		if para.Type != nil {
			f.printExpr(para.Type)
		}
	}
	f.printToken(lst.Rparen)
}

func printFuncSig(f *formatter, fsig *ast.FuncSig) {
	printParaList(f, fsig.Args)
	if fsig.RetType != nil {
		f.printExprs(" ", fsig.RetType)
	} else if fsig.Rets != nil {
		f.printSpace()
		printParaList(f, fsig.Rets)
	}
}

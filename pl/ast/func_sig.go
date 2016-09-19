package ast

import (
	"shanhu.io/smlvm/lexing"
)

// Para is a function parameter
type Para struct {
	Ident *lexing.Token
	Type  Expr // when Type is missing, Ident also might be the type
}

// ParaList is a parameter list
type ParaList struct {
	Lparen *lexing.Token
	Paras  []*Para
	Commas []*lexing.Token
	Rparen *lexing.Token
}

// Named checks if the parameter list is a named list
// or anonymous list.
func (lst *ParaList) Named() bool {
	for _, p := range lst.Paras {
		if p.Ident == nil || p.Type == nil {
			continue
		}
		return true
	}
	return false
}

// Len returns the count of parameters
func (lst *ParaList) Len() int { return len(lst.Paras) }

// FuncSig is a function Signature
type FuncSig struct {
	Args    *ParaList
	Rets    *ParaList // ret list
	RetType Expr      // single ret type
}

// NamedRet returns if the function has a named return
// parameter list
func (sig *FuncSig) NamedRet() bool {
	if sig.Rets == nil {
		return false
	}
	return sig.Rets.Named()
}

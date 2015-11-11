package ir

import (
	"fmt"
	"io"

	"e8vm.io/e8vm/fmt8"
)

func printOp(p io.Writer, op op) {
	switch op := op.(type) {
	case *comment:
		fmt.Fprintf(p, "// %s\n", op.s)
	case *arithOp:
		if op.a == nil {
			if op.op == "" {
				fmt.Fprintf(p, "%s = %s\n", op.dest, op.b)
			} else if op.op == "0" {
				fmt.Fprintf(p, "%s = 0\n", op.dest)
			} else {
				fmt.Fprintf(p, "%s = %s %s\n", op.dest, op.op, op.b)
			}
		} else {
			fmt.Fprintf(p, "%s = %s %s %s\n",
				op.dest, op.a, op.op, op.b,
			)
		}
	case *callOp:
		var args string
		if op.args != nil {
			args = fmt8.Join(op.args, ",")
		}
		fmt.Fprintf(p, "%s = %s(%s)\n", op.dest, op.f, args)
	default:
		panic(fmt.Errorf("invalid or unknown IR op: %T", op))
	}
}

func printJump(p io.Writer, j *blockJump) {
	if j == nil {
		return
	}

	switch j.typ {
	case jmpAlways:
		fmt.Fprintf(p, "goto %s\n", j.to)
	case jmpIf:
		fmt.Fprintf(p, "if %s goto %s\n", j.cond, j.to)
	case jmpIfNot:
		fmt.Fprintf(p, "if !%s goto %s\n", j.cond, j.to)
	default:
		panic("invalid jump type")
	}
}

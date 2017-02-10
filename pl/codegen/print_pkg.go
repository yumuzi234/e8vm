package codegen

import (
	"fmt"
	"io"

	"shanhu.io/smlvm/dasm"
	"shanhu.io/smlvm/fmtutil"
	"shanhu.io/smlvm/lexing"
)

func printBlock(p *fmtutil.Printer, b *Block) {
	fmt.Fprintf(p, "%s:\n", b)
	p.Tab()
	for _, op := range b.ops {
		printOp(p, op)
	}
	printJump(p, b.jump)
	p.ShiftTab()
}

func printBlockInsts(p *fmtutil.Printer, b *Block) {
	fmt.Fprintf(p, "%s:\n", b)
	p.Tab()
	for _, inst := range b.insts {
		fmt.Fprintln(p, dasm.LineStr(inst.inst))
	}
	p.ShiftTab()
}

func printVars(p *fmtutil.Printer, seg string, vars []*Var) {
	if len(vars) == 0 {
		return
	}

	fmt.Fprintf(p, "[%s]\n", seg)
	for _, v := range vars {
		if v.ViaReg == 0 {
			fmt.Fprintf(p, "%s size=%d\n",
				v, v.size,
			)
		} else {
			fmt.Fprintf(p, "%s size=%d reg=%d\n",
				v, v.size, v.ViaReg,
			)
		}
	}
}

func printFunc(p *fmtutil.Printer, f *Func) {
	fmt.Fprintf(p, "func %s {\n", f.name)
	p.Tab()

	printVars(p, "args", f.sig.args)
	printVars(p, "rets", f.sig.rets)
	printVars(p, "saved regs", f.savedRegs)
	printVars(p, "locals", f.locals)

	for b := f.prologue.next; b != f.epilogue; b = b.next {
		printBlock(p, b)
	}

	fmt.Fprintln(p, "----")

	for b := f.prologue.next; b != f.epilogue; b = b.next {
		printBlockInsts(p, b)
	}

	p.ShiftTab()
	fmt.Fprintln(p, "}")
}

func printPkg(p *fmtutil.Printer, pkg *Pkg) {
	fmt.Fprintf(p, "package %s\n", pkg.path)

	for _, f := range pkg.funcs {
		fmt.Fprintln(p)
		printFunc(p, f)
	}
}

// PrintPkg prints a the content of a IR package
func PrintPkg(out io.Writer, pkg *Pkg) error {
	p := fmtutil.NewPrinter(out)
	printPkg(p, pkg)
	return p.Err()
}

// AddDebug adds debug symbols via the add function.
func AddDebug(p *Pkg, add func(f string, pos *lexing.Pos, frame uint32)) {
	for _, f := range p.funcs {
		add(f.name, f.pos, uint32(f.frameSize))
	}
}

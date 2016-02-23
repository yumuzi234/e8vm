package main

import (
	"bytes"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"os"

	"e8vm.io/e8vm/arch8"
	"e8vm.io/e8vm/dasm8"
	"e8vm.io/e8vm/g8"
	"e8vm.io/e8vm/g8/gfmt"
	"e8vm.io/e8vm/g8/parse"
	"e8vm.io/e8vm/lex8"
)

func exit(e error) {
	if e != nil {
		fmt.Fprintln(os.Stderr, e)
	}
	os.Exit(-1)
}

func printErrs(es []*lex8.Error) {
	if len(es) == 0 {
		return
	}
	for _, e := range es {
		fmt.Println(e)
	}
	exit(nil)
}

var (
	bare     = flag.Bool("bare", false, "parse as bare function")
	parseAST = flag.Bool("parse", false, "parse only and print out the ast")
	ir       = flag.Bool("ir", false, "prints out the IR")
	dasm     = flag.Bool("d", false, "deassemble the image")
	ncycle   = flag.Int("n", 100000, "maximum number of cycles")
	verbose  = flag.Bool("v", false, "verbose")
	golike   = flag.Bool("golike", false, "using strict go-like syntax")
)

func main() {
	flag.Parse()

	args := flag.Args()
	if len(args) != 1 {
		exit(errors.New("need exactly one input input file"))
	}
	fname := args[0]
	input, e := ioutil.ReadFile(fname)
	if e != nil {
		exit(e)
	}

	if *bare {
		if *parseAST {
			stmts, es := parse.Stmts(fname, bytes.NewBuffer(input))
			printErrs(es)
			gfmt.FprintStmts(os.Stdout, stmts)
		} else {
			bs, es, irLog := g8.CompileBareFunc(fname, string(input))
			printErrs(es)
			printIRLog(irLog, *ir)
			runImage(bs, *dasm, *ncycle)
		}
	} else {
		if *parseAST {
			f, rec, es := parse.File(fname, bytes.NewBuffer(input), true)
			printErrs(es)
			gfmt.FprintFile(os.Stdout, f, rec)
		} else {
			bs, es, irLog := g8.CompileAndTestSingle(
				fname, string(input), *golike,
			)
			printErrs(es)
			printIRLog(irLog, *ir)
			runImage(bs, *dasm, *ncycle)
		}
	}
}

func runImage(bs []byte, dasm bool, n int) {
	if dasm {
		err := dasm8.DumpImage(bytes.NewReader(bs), os.Stdout)
		if err != nil {
			fmt.Println(err)
		}
	}
	if len(bs) == 0 {
		fmt.Println("(the image is empty)")
		return
	}

	m := arch8.NewMachine(0, 1)
	if err := m.LoadImageBytes(bs); err != nil {
		fmt.Println(err)
		return
	}

	ncycle, exp := m.Run(n)
	fmt.Printf("(%d cycles)\n", ncycle)
	if exp != nil {
		if !arch8.IsHalt(exp) {
			fmt.Println(exp)
			err := arch8.FprintStack(os.Stdout, m, exp)
			if err != nil {
				fmt.Println(err)
			}
		}
	} else {
		fmt.Println("(end of time)")
	}
}

func printIRLog(irLog []byte, ir bool) {
	if !ir {
		return
	}
	if irLog == nil {
		fmt.Println("(no IR log produced)")
	} else {
		fmt.Println(string(irLog))
	}
}

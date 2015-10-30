package main

import (
	"flag"
	"fmt"
	"log"
	"math"
	"os"
	"runtime/pprof"

	"e8vm.io/e8vm/arch8"
	"e8vm.io/e8vm/asm8"
	"e8vm.io/e8vm/build8"
	"e8vm.io/e8vm/g8"
)

var (
	golike = flag.Bool("golike", false, "uses go-like syntax")
	doTest = flag.Bool("test", true, "also run tests")
	initPC = flag.Uint("initpc", arch8.InitPC,
		"the starting address of the image",
	)
	cpuProfile = flag.String("profile", "", "cpu profile output")
)

func checkInitPC() {
	if *initPC > math.MaxUint32 {
		fmt.Fprintln(os.Stderr, "init pc out of range")
		os.Exit(-1)
	}
	if *initPC%arch8.RegSize != 0 {
		fmt.Fprintln(os.Stderr, "init pc not aligned")
		os.Exit(-1)
	}
}

func main() {
	flag.Parse()
	if *cpuProfile != "" {
		f, err := os.Create(*cpuProfile)
		if err != nil {
			log.Fatal(err)
		}
		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
	}

	checkInitPC()

	var lang build8.Lang
	if !*golike {
		lang = g8.Lang()
	} else {
		lang = g8.LangGolike()
	}

	home := build8.NewDirHome(".", lang)
	home.AddLang("asm", asm8.Lang())
	home.AddLang("bare", g8.BareFunc())

	b := build8.NewBuilder(home)
	b.Verbose = true
	b.InitPC = uint32(*initPC)

	es := b.BuildAll(*doTest)
	if es != nil {
		for _, e := range es {
			fmt.Println(e)
		}
		os.Exit(-1)
	}
}

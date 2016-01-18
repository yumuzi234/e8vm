package main

import (
	"flag"
	"fmt"
	"log"
	"math"
	"os"
	"runtime/pprof"
	"strings"

	"e8vm.io/e8vm/arch8"
	"e8vm.io/e8vm/asm8"
	"e8vm.io/e8vm/build8"
	"e8vm.io/e8vm/g8"
	"e8vm.io/e8vm/lex8"
)

var (
	golike     = flag.Bool("golike", false, "uses go-like syntax")
	runTests   = flag.Bool("test", true, "also run tests")
	staticOnly = flag.Bool("static", false, "do static analysis only")
	initPC     = flag.Uint("initpc", arch8.InitPC,
		"the starting address of the image",
	)
	cpuProfile = flag.String("profile", "", "cpu profile output")
	pkg        = flag.String("pkg", "", "package to build")
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
		lang = g8.LangGoLike()
	}

	home := build8.NewDirHome(".", lang)
	home.AddLang("asm", asm8.Lang())
	home.AddLang("bare", g8.BareFunc())

	b := build8.NewBuilder(home, home)
	b.Verbose = true
	b.InitPC = uint32(*initPC)
	b.RunTests = *runTests
	b.StaticOnly = *staticOnly

	var es []*lex8.Error
	if *pkg == "" {
		es = b.BuildAll()
	} else if strings.HasSuffix(*pkg, "...") {
		prefix := strings.TrimSuffix(*pkg, "...")
		b.BuildPrefix(prefix)
	} else {
		es = b.Build(*pkg)
	}

	if es != nil {
		for _, e := range es {
			fmt.Println(e)
		}
		os.Exit(-1)
	}
}

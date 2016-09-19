package main

import (
	"flag"
	"fmt"
	"log"
	"math"
	"os"
	"runtime/pprof"
	"strings"

	"shanhu.io/smlvm/arch"
	"shanhu.io/smlvm/asm"
	"shanhu.io/smlvm/builds"
	"shanhu.io/smlvm/lexing"
	"shanhu.io/smlvm/pl"
)

var (
	golike     = flag.Bool("golike", false, "uses go-like syntax")
	runTests   = flag.Bool("test", true, "also run tests")
	staticOnly = flag.Bool("static", false, "do static analysis only")
	initPC     = flag.Uint("initpc", arch.InitPC,
		"the starting address of the image",
	)
	cpuProfile = flag.String("profile", "", "cpu profile output")
	pkg        = flag.String("pkg", "", "package to build")
	homeDir    = flag.String("home", ".", "the home directory")
)

func checkInitPC() {
	if *initPC > math.MaxUint32 {
		fmt.Fprintln(os.Stderr, "init pc out of range")
		os.Exit(-1)
	}
	if *initPC%arch.RegSize != 0 {
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

	lang := pl.Lang(*golike)
	home := builds.NewDirHome(*homeDir, lang)
	home.AddLang("asm", asm.Lang())
	home.AddLang("bare", pl.BareFunc())

	b := builds.NewBuilder(home, home)
	b.Verbose = true
	b.InitPC = uint32(*initPC)
	b.RunTests = *runTests
	b.StaticOnly = *staticOnly

	var es []*lexing.Error
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

package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"e8vm.io/e8vm/arch8"
	"e8vm.io/e8vm/build8"
	"e8vm.io/e8vm/home8"
	"e8vm.io/e8vm/lex8"
)

var (
	runTests = flag.Bool("test", true, "run tests")
	pkg      = flag.String("pkg", "", "package to build")
	homeDir  = flag.String("home", ".", "the home directory")
)

func main() {
	flag.Parse()

	home := home8.NewHome(*homeDir, "")
	b := build8.NewBuilder(home, home)
	b.Verbose = true
	b.InitPC = arch8.InitPC
	b.RunTests = *runTests

	var es []*lex8.Error
	if *pkg == "" {
		es = b.BuildPrefix("/")
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

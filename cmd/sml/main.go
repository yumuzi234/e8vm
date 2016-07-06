package main

import (
	"flag"
	"fmt"
	"os"

	"e8vm.io/e8vm/arch8"
	"e8vm.io/e8vm/build8"
	"e8vm.io/e8vm/home8"
	"e8vm.io/e8vm/lex8"
)

var (
	runTests = flag.Bool("test", true, "run tests")
	pkg      = flag.String("pkg", "/...", "package to build")
	homeDir  = flag.String("home", ".", "the home directory")
	plan     = flag.Bool("plan", false, "plan only")
)

func handleErrs(errs []*lex8.Error) {
	if errs == nil {
		return
	}
	for _, err := range errs {
		fmt.Println(err)
	}
	os.Exit(-1)
}

func main() {
	flag.Parse()

	home := home8.NewHome(*homeDir, "")
	b := build8.NewBuilder(home, home)
	b.Verbose = true
	b.InitPC = arch8.InitPC
	b.RunTests = *runTests

	pkgs, err := build8.SelectPkgs(home, *pkg)
	if err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}

	if !*plan {
		handleErrs(b.BuildPkgs(pkgs))
	} else {
		buildOrder, errs := b.Plan(pkgs)
		handleErrs(errs)
		for _, p := range buildOrder {
			fmt.Println(p)
		}
	}
}

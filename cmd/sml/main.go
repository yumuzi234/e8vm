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
	plan     = flag.Bool("plan", false, "plan only")
)

func selectPkgs(h *home8.Home, s string) ([]string, error) {
	if s == "" {
		return h.Pkgs("/"), nil
	}
	if strings.HasSuffix(s, "...") {
		prefix := strings.TrimSuffix(s, "...")
		return h.Pkgs(prefix), nil
	}

	if !h.HasPkg(s) {
		return nil, fmt.Errorf("package %q not found", s)
	}
	return []string{s}, nil
}

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

	pkgs, err := selectPkgs(home, *pkg)
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

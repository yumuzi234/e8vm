package build8

import (
	"fmt"
	"sort"
	"strings"

	"e8vm.io/e8vm/arch8"
	"e8vm.io/e8vm/lex8"
)

func runTests(
	log lex8.Logger, tests map[string]uint32, img []byte, verbose bool,
) {
	report := func(name string, pass bool, err error) {
		if !pass {
			lex8.LogError(log, fmt.Errorf("%s failed: got %s", name, err))
			if verbose {
				fmt.Println("FAILED")
			}
			return
		}

		if verbose {
			fmt.Println("pass")
		}
	}

	var testNames []string
	for name := range tests {
		testNames = append(testNames, name)
	}
	sort.Strings(testNames)

	for _, test := range testNames {
		if verbose {
			fmt.Printf("  - %s: ", test)
		}

		arg := tests[test]
		_, err := arch8.RunImageArg(img, arg)
		if strings.HasPrefix(test, "TestBad") {
			report(test, arch8.IsPanic(err), err)
		} else {
			report(test, arch8.IsHalt(err), err)
		}
	}
}

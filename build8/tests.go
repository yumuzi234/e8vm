package build8

import (
	"errors"
	"fmt"
	"sort"
	"strings"

	"e8vm.io/e8vm/arch8"
	"e8vm.io/e8vm/lex8"
)

var errTimeOut = errors.New("time out")

func cycleStr(n int) string {
	if n >= -1 && n <= 1 {
		return fmt.Sprintf("%d cycle", n)
	}
	return fmt.Sprintf("%d cycles", n)
}

func runTests(
	log lex8.Logger, tests map[string]uint32, img []byte,
	verbose bool, ncycle int, logln func(s string),
) {
	// TODO(h8liu): this reporting should go with JSON for better formatting.
	report := func(name string, ncycle int, pass bool, err error) {
		if !pass {
			if err == nil {
				err = errTimeOut
			}
			lex8.LogError(log, fmt.Errorf("%s failed: got %s", name, err))
			if verbose {
				logln(fmt.Sprintf(
					"  - %s: FAILED (%s, got %s)",
					name, cycleStr(ncycle), err,
				))
			}
			return
		}

		if verbose {
			logln(fmt.Sprintf("  - %s: passed (%s)", name, cycleStr(ncycle)))
		}
	}

	var testNames []string
	for name := range tests {
		testNames = append(testNames, name)
	}
	sort.Strings(testNames)

	for _, test := range testNames {
		arg := tests[test]
		n, err := arch8.RunImageArg(img, arg, ncycle)
		if strings.HasPrefix(test, "TestBad") {
			report(test, n, arch8.IsPanic(err), err)
		} else {
			report(test, n, arch8.IsHalt(err), err)
		}
	}
}

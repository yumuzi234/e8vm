package build8

import (
	"bytes"
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
	report := func(
		name string, ncycle int, pass bool,
		m *arch8.Machine, err error,
	) {
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
				// TODO(h8liu): this is too ugly here...
				excep, ok := err.(*arch8.CoreExcep)
				if !arch8.IsHalt(err) && ok {
					stackTrace := new(bytes.Buffer)
					arch8.FprintStack(stackTrace, m, excep)
					logln(stackTrace.String())
				}
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
		m := arch8.NewMachine(0, 1)
		if err := m.LoadImageBytes(img); err != nil {
			report(test, 0, false, m, err)
			continue
		}
		if err := m.WriteWord(arch8.AddrBootArg, arg); err != nil {
			report(test, 0, false, m, err)
			continue
		}

		var err error
		n, excep := m.Run(ncycle)
		if excep == nil {
			err = errTimeOut
		}
		if strings.HasPrefix(test, "TestBad") {
			report(test, n, arch8.IsPanic(err), m, err)
		} else {
			report(test, n, arch8.IsHalt(err), m, err)
		}
	}
}

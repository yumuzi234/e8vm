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
	log lex8.Logger, tests map[string]uint32, img []byte, opt *Options,
) {
	logln := func(s string) {
		if opt.LogLine == nil {
			fmt.Println(s)
		} else {
			opt.LogLine(s)
		}
	}

	// TODO(h8liu): this reporting should go with JSON for better formatting.
	report := func(
		name string, ncycle int, pass bool,
		m *arch8.Machine, err error,
	) {
		if pass {
			if opt.Verbose {
				logln(fmt.Sprintf(
					"  - %s: passed (%s)", name, cycleStr(ncycle),
				))
			}
			return
		}

		if err == nil {
			err = errTimeOut
		}
		lex8.LogError(log, fmt.Errorf("%s failed: got %s", name, err))
		if opt.Verbose {
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
	}

	var testNames []string
	for name := range tests {
		testNames = append(testNames, name)
	}
	sort.Strings(testNames)

	for _, test := range testNames {
		arg := tests[test]
		m := arch8.NewMachine(&arch8.Config{
			BootArg: arg,
		})
		if err := m.LoadImageBytes(img); err != nil {
			report(test, 0, false, m, err)
			continue
		}

		var err error
		n, excep := m.Run(opt.TestCycles)
		if excep == nil {
			err = errTimeOut
		} else {
			err = excep
		}
		if strings.HasPrefix(test, "TestBad") {
			report(test, n, arch8.IsPanic(err), m, err)
		} else {
			report(test, n, arch8.IsHalt(err), m, err)
		}
	}
}

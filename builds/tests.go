package builds

import (
	"bytes"
	"errors"
	"fmt"
	"sort"
	"strings"

	"shanhu.io/smlvm/arch"
	"shanhu.io/smlvm/lexing"
)

var errTimeOut = errors.New("time out")

func cycleStr(n int) string {
	if n >= -1 && n <= 1 {
		return fmt.Sprintf("%d cycle", n)
	}
	return fmt.Sprintf("%d cycles", n)
}

func runTests(
	log lexing.Logger, tests map[string]uint32, img []byte, opt *Options,
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
		m *arch.Machine, err error,
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
		lexing.LogError(log, fmt.Errorf("%s failed: got %s", name, err))
		if opt.Verbose {
			logln(fmt.Sprintf(
				"  - %s: FAILED (%s, got %s)",
				name, cycleStr(ncycle), err,
			))
			// TODO(h8liu): this is too ugly here...
			excep, ok := err.(*arch.CoreExcep)
			if !arch.IsHalt(err) && ok {
				stackTrace := new(bytes.Buffer)
				arch.FprintStack(stackTrace, m, excep)
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
		m := arch.NewMachine(&arch.Config{
			BootArg: arg,
			InitPC:  opt.InitPC,
			InitSP:  opt.InitSP,
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
			report(test, n, arch.IsPanic(err), m, err)
		} else {
			report(test, n, arch.IsHalt(err), m, err)
		}
	}
}

func runPkgTests(c *context, p *pkg) []*lexing.Error {
	lib := p.pkg.Lib
	tests := p.pkg.Tests
	testMain := p.pkg.TestMain
	if testMain != "" && lib.HasFunc(testMain) {
		log := lexing.NewErrorList()
		if len(tests) > 0 {
			bs := new(bytes.Buffer)
			lexing.LogError(log, linkPkg(c, bs, p, testMain))
			fout := c.output.TestBin(p.path)

			img := bs.Bytes()
			_, err := fout.Write(img)
			lexing.LogError(log, err)
			lexing.LogError(log, fout.Close())
			if es := log.Errs(); es != nil {
				return es
			}

			runTests(log, tests, img, c.Options)
			if es := log.Errs(); es != nil {
				return es
			}
		}
	}

	return nil
}

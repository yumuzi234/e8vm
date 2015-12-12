// gfmt is the code formatter of G language.
package main

import (
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"e8vm.io/e8vm/g8/gfmt"
	"e8vm.io/e8vm/g8/parse"
)

const tabExpand = 4

func countIndent(line string) int {
	ret := 0

loop:
	for _, r := range line {
		switch r {
		case ' ':
			ret++
		case '\t':
			ret += tabExpand
		default:
			break loop
		}
	}

	return ret
}

func fmtLine(line string) string {
	ret := new(bytes.Buffer)

	indent := countIndent(line)
	for i := 0; i < indent; i++ {
		ret.WriteRune(' ')
	}

	ret.WriteString(strings.TrimLeft(line, " \t"))

	retLine := ret.String()
	return retLine
}

var tempDir = os.TempDir()

func fmtFile(fname string) (bool, error) {
	input, e := ioutil.ReadFile(fname)
	if e != nil {
		return false, e
	}

	f, rec, es := parse.File(fname, bytes.NewBuffer(input), false)
	if es != nil {
		return false, fmt.Errorf("%d errors found at parsing", len(es))
	}

	var output bytes.Buffer
	gfmt.FprintFile(&output, f, rec)
	if bytes.Compare(input, output.Bytes()) == 0 {
		return false, nil
	}

	tempfile, e := ioutil.TempFile(tempDir, "e8fmt")
	if e != nil {
		return false, e
	}
	if _, e := tempfile.Write(output.Bytes()); e != nil {
		return false, e
	}
	if e := tempfile.Close(); e != nil {
		return false, e
	}
	return true, os.Rename(tempfile.Name(), fname)
}

func main() {
	flag.Parse()

	args := flag.Args()
	for _, fname := range args {
		if changed, e := fmtFile(fname); e != nil {
			log.Print(e)
		} else if changed {
			fmt.Println(fname)
		}
	}
}

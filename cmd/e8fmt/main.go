// e8fmt is the code formatter of G language.
package main

import (
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"e8vm.io/e8vm/g8/ast"
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

	f, es := parse.File(fname, bytes.NewBuffer(input), false)
	if es != nil {
		return false, fmt.Errorf("%d errors found at parsing", len(es))
	}

	temp, e := ioutil.TempFile(tempDir, "e8fmt")
	if e != nil {
		return false, e
	}

	ast.FprintFile(temp, f)

	temp.Seek(0, os.SEEK_SET)
	output, e := ioutil.ReadAll(temp)
	if e != nil {
		return false, e
	}
	changed := bytes.Compare(input, output) != 0

	if e := temp.Close(); e != nil {
		return false, e
	}

	if !changed {
		return false, nil
	}
	return true, os.Rename(temp.Name(), fname)
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

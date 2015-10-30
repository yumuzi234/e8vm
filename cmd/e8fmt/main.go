// e8fmt is the code formatter of G language.
package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
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
	f, e := os.Open(fname)
	if e != nil {
		return false, e
	}

	temp, e := ioutil.TempFile(tempDir, "e8fmt")
	if e != nil {
		return false, e
	}

	needMove := false
	s := bufio.NewScanner(f)
	for s.Scan() {
		line := s.Text()
		formatted := fmtLine(line)
		fmt.Fprintln(temp, formatted)
		if line != formatted {
			needMove = true
		}
	}

	if e := s.Err(); e != nil {
		return false, e
	}
	if e := f.Close(); e != nil {
		return false, e
	}
	if e := temp.Close(); e != nil {
		return false, e
	}

	if !needMove {
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

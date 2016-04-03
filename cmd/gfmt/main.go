// gfmt is the code formatter of G language.
package main

import (
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"e8vm.io/e8vm/g8/gfmt"
)

var tempDir = os.TempDir()

func fmtFile(fname string) (bool, error) {
	input, e := ioutil.ReadFile(fname)
	if e != nil {
		return false, e
	}

	out, errs := gfmt.File(fname, input)
	if errs != nil {
		for _, err := range errs {
			fmt.Println(err)
		}
		return false, fmt.Errorf("%d errors found at parsing", len(errs))
	}

	if bytes.Compare(input, out) == 0 {
		return false, nil
	}

	tempfile, e := ioutil.TempFile(tempDir, "e8fmt")
	if e != nil {
		return false, e
	}
	if _, e := tempfile.Write(out); e != nil {
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

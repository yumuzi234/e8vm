// gfmt is the code formatter of G language.
package main

import (
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"

	"shanhu.io/smlvm/builds"
	"shanhu.io/smlvm/pl"
	"shanhu.io/smlvm/pl/gfmt"
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

	tempfile, err := ioutil.TempFile(tempDir, "gfmt")
	if err != nil {
		return false, err
	}
	if _, err := tempfile.Write(out); err != nil {
		return false, err
	}
	if err := tempfile.Close(); err != nil {
		return false, err
	}
	return true, os.Rename(tempfile.Name(), fname)
}

func main() {
	flag.Parse()

	args := flag.Args()
	if len(args) == 0 {
		in := builds.NewDirFS(".")
		langSet := pl.MakeLangSet(false)
		pkgs, err := builds.SelectPkgs(in, langSet, "")
		if err != nil {
			log.Print(err)
			return
		}

		for _, pkg := range pkgs {
			p := strings.TrimPrefix(pkg, "/")
			files, err := builds.ListSrcFiles(in, langSet, p)
			if err != nil {
				log.Print(err)
				continue
			}

			for _, file := range files {
				name := filepath.FromSlash(path.Join(p, file))
				changed, err := fmtFile(name)
				if err != nil {
					log.Print(err)
				} else if changed {
					fmt.Println(name)
				}
			}
		}
	} else {
		for _, fname := range args {
			if changed, err := fmtFile(fname); err != nil {
				log.Print(err)
			} else if changed {
				fmt.Println(fname)
			}
		}
	}
}

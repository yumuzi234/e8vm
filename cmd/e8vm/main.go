package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"os"

	"e8vm.io/e8vm/arch8"
	"e8vm.io/e8vm/dasm8"
	"e8vm.io/e8vm/debug8"
	"e8vm.io/e8vm/e8"
)

var (
	doDasm      = flag.Bool("d", false, "do dump")
	printDebug  = flag.Bool("debug", false, "print debug symbols")
	ncycle      = flag.Int("n", 100000, "max cycles to execute")
	memSize     = flag.Int("m", 0, "memory size; 0 for full 4GB")
	printStatus = flag.Bool("s", false, "print status after execution")
	bootArg     = flag.Uint("arg", 0, "boot argument, a uint32 number")
	romRoot     = flag.String("rom", "", "rom root path")
	randSeed    = flag.Int64("seed", 0, "random seed, 0 for using the time")
)

func run(bs []byte) (int, error) {
	// create a single core machine
	m := arch8.NewMachine(uint32(*memSize), 1)
	if err := m.LoadImageBytes(bs); err != nil {
		return 0, err
	}

	if *bootArg > math.MaxUint32 {
		log.Fatalf("boot arg(%d) is too large", *bootArg)
	}
	if err := m.WriteWord(arch8.AddrBootArg, uint32(*bootArg)); err != nil {
		return 0, err
	}

	if *romRoot != "" {
		m.MountROM(*romRoot)
	}
	if *randSeed != 0 {
		m.RandSeed(*randSeed)
	}

	ret, exp := m.Run(*ncycle)
	if *printStatus {
		m.PrintCoreStatus()
	}

	if exp == nil {
		return ret, nil
	}
	return ret, exp
}

func main() {
	flag.Parse()

	args := flag.Args()
	if len(args) != 1 {
		log.Fatal("need exactly one input file\n")
	}

	fname := args[0]

	if *doDasm {
		f, err := os.Open(fname)
		defer f.Close()

		err = dasm8.DumpImage(f, os.Stdout)
		if err != nil {
			log.Fatal(err)
		}
	} else if *printDebug {
		f, err := os.Open(fname)
		defer f.Close()

		secs, err := e8.Read(f)
		if err != nil {
			log.Fatal(err)
		}

		for _, sec := range secs {
			if sec.Type != e8.Debug {
				continue
			}

			tab, err := debug8.UnmarshalTable(sec.Bytes)
			if err != nil {
				log.Fatal(err)
			}

			tab.PrintTo(os.Stdout)
		}
	} else {
		bs, err := ioutil.ReadFile(fname)
		if err != nil {
			log.Fatal(err)
		}
		n, e := run(bs)
		fmt.Printf("(%d cycles)\n", n)
		if e != nil {
			if !arch8.IsHalt(e) {
				fmt.Println(e)
			}
		} else {
			fmt.Println("(end of time)")
		}
	}
}

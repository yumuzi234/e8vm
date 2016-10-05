[![BuildStatus](https://travis-ci.org/shanhuio/smlvm.png?branch=master)](https://travis-ci.org/shanhuio/smlvm)

```
go get -u shanhu.io/smlvm/...
```

# Small VM

Small Virtual Machine (smlvm) is a self-contained system that has its
own instruction set, assembly language and assembler, system language,
and project building system.

The project is written entirely in Go language. Each file in the
project has no more than 300 lines, with each line no more than 80
characters. Among these small files, there are no circular
dependencies, checked by static analysis, and as a result, the project
architecture can be automatically visualized from static code
analysis.

The main project in this repository depends on nothing other than the
Go standard library. It is *NOT* yet another compiler project based on
LLVM.

For Go language documentation on the package APIs, I recommend
[GoWalker](https://gowalker.org/shanhu.io/smlvm). I find it slightly
better than [godoc.org](https://godoc.org/shanhu.io/smlvm).

## To Use `make`

The project comes with a `makefile`, which formats the code files,
check lints, check circular dependencies and build tags. Running the
`makefile` requires installing some tools.

```
go get -u shanhu.io/smlvm/...
go get -u github.com/golang/lint/golint
go get -u github.com/jstemmer/gotags
```

## Copyright and License

The project developers own the copyright; my employer (Google) does
*NOT* own the copyright. The Licence is Apache.

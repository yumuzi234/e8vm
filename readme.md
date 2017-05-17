[![BuildStatus](https://travis-ci.org/shanhuio/smlvm.png?branch=master)](https://travis-ci.org/shanhuio/smlvm)

# Small VM

Small Virtual Machine (`smlvm`) provides a simple programming language
that compiles to a simulated, simple virtual machine. It is
essentially a subset of Go programming language.

[Try in playground](https://smallrepo.com/play)

## Install

```
go get -u shanhu.io/smlvm/...
```

## Caveats

This project is written in Go. Each source file in this project has no
more than 300 lines (80 max per line). Also there are no circular
dependencies among files. As a result, the project architecture can be
[visualized](https://shanhu.io/smlvm).

This repository does not use any third party libraries; it depends on
only Go standard library. The compiler is all written from scratch,
and does not use LLVM.

[Package Docs](https://godoc.org/shanhu.io/smlvm).

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

Copyright by Shanhu contributors. Licence Apache.

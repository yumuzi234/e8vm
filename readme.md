[![BuildStatus](https://travis-ci.org/e8vm/e8vm.png?branch=master)](https://travis-ci.org/h8liu/e8vm)

```
go get -u e8vm.io/e8vm/...
```

# E8VM

E8VM stands for Emul8ed Virtual Machine. It is a self-contained system
that has its own instruction set -- `arch8`, its assembly language and
assembler -- `asm8`, its own system language -- `g8`, and its own
project building system -- `build8`.

The project is written entirely in Go language. Plus, each file in the
project has no more than 300 lines, with each line no more than 80
characters. Among these small files, there are no circular
dependencies, and as a result, the project architecture can be
automatically [visualized](http://8k.lonnie.io) from static code
analysis.

For Go language documentation on the package APIs, I recommend
[GoWalker](https://gowalker.org/e8vm.io/e8vm). I find it slightly
better than [godoc.org](https://godoc.org/e8vm.io/e8vm).

## To Use `make`

The project comes with a `makefile`, which formats the code files,
check lints, check circular dependencies and build tags. Running the
`makefile` requires installing some tools.

```
go get -u e8vm.io/tools/...
go get -u github.com/golang/lint/golint
go get -u github.com/jstemmer/gotags
```

## Copyright and License

The project developers own the copyright; my empolyer (Google) does
not own the copyright. Apache is the License.

[![BuildStatus](https://travis-ci.org/shanhuio/smlvm.png?branch=master)](https://travis-ci.org/shanhuio/smlvm)

```
go get -u shanhu.io/smlvm/...
```

# Small VM and G Language

Small Virtual Machine (`smlvm`) is a self-contained system that has
its own instruction set, assembly language and assembler, system
language, and project building system.

The main project in this repository depends on nothing other than the
Go standard library. It is *NOT* yet another compiler project based on
LLVM.

[Introduction of G]

For Go language documentation on the package APIs, I recommend
[GoWalker](https://gowalker.org/shanhu.io/smlvm). I find it slightly
better than [godoc.org](https://godoc.org/shanhu.io/smlvm).

## Why Small VM

### Approachable Compiler and Evolving Language

We build Small VM because

The project is written entirely in Go language. To make the project is 
clean and approachable, each file in the project has no more than 300 lines, 
with each line no more than characters. Among these small files, 
there are no circular dependencies, checked by static analysis, 
and as a result, the project architecture can be automatically
[visualized](https://shanhu.io/smlvm) from static code analysis.

### Cross Platform?
smlrepo?


## To Use 

### `make`

The project comes with a `makefile`, which formats the code files,
check lints, check circular dependencies and build tags. Running the
`makefile` requires installing some tools.

```
go get -u shanhu.io/smlvm/...
go get -u github.com/golang/lint/golint
go get -u github.com/jstemmer/gotags
```
### playground


## Copyright and License

Copyright by Shanhu contributors. Licence Apache.

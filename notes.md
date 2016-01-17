## Static analysis

- we need to support building a particular set of package
- like a single package or all packages with a particular prefix
- we also need to support saving static analysis results
- there should be a common syntax saving format for both golang and g language
- eventually, there should also be a search engine

- our build system is a little bit messy
- we are not clear about what it the input and what is the output 

- clearly the input could be a file system
- but it also could be a virtual file system
- it is fundamentally a set of packages, where each package has a set of files
- source files, or even non source files
- we can have a package selector
- like one single package, by a name
- or a list of packages that has the common prefix
- or all packages in a repository

- how about the output
- the output are files of syntaxed parsted tokens
- identifer tokens are either anchors
- or links to another identifer
- all global identifers has a global name for routing
- local identifiers also has a number
- for each packages that are successfully compiled, it can output a library
- it can also save itself for caching the output result, but this is optional
- each package also have a bunch of tests that can be run
- each test case also has a name
- test cases will be runned
- and finally a package might produce a binary

- so the output has

input: a list of packages
options: build depths
for each packages
- reads in files
- or read in saved libs
- build errors if any (to stdout)
- parsed files (to parsed)
- dependency structure (to deps)
- symbol/identifier index (to sym)
- test results (also to stdout, but could be structured)
- built binaries (to bin)

bin: built binaries
pkg: log, pkg, deps, depmap, syms

build

SrcDir {
	List(ext string) []string
	Open(name string) (io.Reader, error)
}

Output {
	
}
Binary io.Writer

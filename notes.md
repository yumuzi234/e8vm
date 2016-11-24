what is a build system?
- a build system should be a computational graph.
- the compiler should be an external computational service.
- a package is a build unit, a node in the graph.
- but for this to work, we first need to build the graph.
- that is, we need to pre-parse the files for the imports.
- then we can register a node with its dependencies.
- source dependencies.
- compiling a package actually does not dependend on the import's source
- it depends on the input's compiling result
- how about tokens and logs
- everything should be immutable
- that is, everything can only be written once.

- so, there is another layer before actually building
- which is the import parsing, the file system layer
- the file system is immutable, where everything can only be written once.
- after it is written, it can be opened and read.
- do we support listing??
- we need to support iterating over all the packages under a particular prefix
- and this iterating, is really just over the current repo, kind of.
- so, repo name should really be an alias,
- where real repos are random ids.
- every new workspace will automatically create an anonymous repo
- which has an id, but does not have an alias
- a workspace without a repo name can commit, but cannot be imported

- what we need is to have the repo concept builtin into our build system

--

- so, instead of wrapping a generic network interface
- let's provide a generic RPC interface
- where it is sync RPC interface by default
- and in the future, it can be just an async RPC interface of some sort.
- an RPC interface
- flag, reserve, service.method. payload addr, payload length
  return addr, return size

0: flag
1: error code
4-6: service id
6-8: method id
8-12: request addr
12-16: request length
16-20: response addr
20-24: response length
24-28: response length (used)
28-32: reserved

--

- we need a token marking
- also for each token

## Static analysis

- we need to support building a particular set of package
- like a single package or all packages with a particular prefix
- we also need to support saving static analysis results

The result of static analysis:

- syntax annotated lines, with identifier links

- file level dependency
- token of each file
- a database that has type and reference of each identifier
- defined identity in each file: const, struct, func and var
- public interface for each package, go doc like sorted
- package import

```
struct file {
	path string
	name string
	depName string

	items []*item
	defines []*ident
	refs map[*ident]*ident
}

struct package {
	path string
	name string

	imports []*import
	files []*file
}

struct ident {
	pos *pos
	refs []*ident
}

```


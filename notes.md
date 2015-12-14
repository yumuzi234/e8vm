## Panic Stack trace

to recover stack trace of at least function calls, we need to save each function:
- package and symbol name
- start pc
- code size
- frame size

and for each function call:
- calling identifier that generate function call

## Static analysis

- what we want is the type of each identifier.
- because of type infering, this means the type of each expression.

The result of static analysis:

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

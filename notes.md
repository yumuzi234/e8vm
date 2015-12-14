## Init support

One thing that is hard about init support is that.

- should init be a G language concept or a general building concept?
- if init is a general concept, then the builder needs to know how to
  build a function entrance.
- which is not really right.
- G language builds on top of asm8.
- if we need a builder that would support other languages upon asm8
- then each language would have its own way to init its runtime when building
- a main procedure.

so, building main is a language function, but not a package function?

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

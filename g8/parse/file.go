package parse

import (
	"bytes"
	"io"
	"io/ioutil"

	"e8vm.io/e8vm/g8/ast"
	"e8vm.io/e8vm/lexing"
	"e8vm.io/e8vm/textbox"
)

func parseTopDecl(p *parser) ast.Decl {
	if p.SeeKeyword("const") {
		return parseConstDecls(p)
	} else if p.SeeKeyword("var") {
		return parseVarDecls(p)
	} else if p.SeeKeyword("func") {
		return parseFunc(p)
	} else if !p.golike && p.SeeKeyword("struct") {
		return parseStruct(p)
	} else if !p.golike && p.SeeLit(Ident, "type") {
		return parseStruct(p)
	} else if p.SeeKeyword("type") && p.golike {
		return parseStruct(p)
	} else if p.SeeKeyword("import") {
		p.ErrorfHere("only one import block allowed at the head")
	}

	if len(p.Errs()) == 0 {
		// we only complain about this when there is no other error
		p.ErrorfHere("expect top level declaration")
	} else {
		p.Jail()
	}
	p.Next() // make some progress anyway
	return nil
}

func parseFile(p *parser) *ast.File {
	ret := &ast.File{Path: p.f}

	if p.golike {
		kw := p.ExpectKeyword("package")
		name := p.Expect(Ident)
		semi := p.ExpectSemi()
		if p.InError() {
			return ret
		}

		ret.Package = &ast.PackageTitle{Kw: kw, Name: name, Semi: semi}
	} else if p.SeeLit(Ident, "package") {
		p.ErrorfHere("G language does not need the package clause")
		p.ErrorfHere("please just remove it")
		p.Next()
		p.Expect(Ident)
		p.ExpectSemi()
		if p.InError() {
			return ret
		}
	}

	if p.SeeKeyword("import") {
		ret.Imports = parseImports(p)
	}

	for !p.See(lexing.EOF) {
		decl := parseTopDecl(p)
		if decl != nil {
			ret.Decls = append(ret.Decls, decl)
		}

		if p.InError() {
			p.skipErrStmt()
		}
	}

	return ret
}

const (
	// MaxLine is the max number of lines of a G language file.
	MaxLine = 300

	// MaxCol is the max number of columns of a G language file.
	MaxCol = 80
)

// File parses a file.
func File(f string, r io.Reader, golike bool) (
	*ast.File, *lexing.Recorder, []*lexing.Error,
) {
	bs, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, nil, lexing.SingleErr(err)
	}

	es := textbox.CheckRect(f, bytes.NewReader(bs), MaxLine, MaxCol)
	if es != nil {
		return nil, nil, es
	}

	p, rec := makeParser(f, bytes.NewReader(bs), golike)
	ret := parseFile(p)
	if es := p.Errs(); es != nil {
		return nil, nil, es
	}
	return ret, rec, nil
}

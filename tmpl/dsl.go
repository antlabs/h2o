package tmpl

import (
	_ "embed"
	"io"
	"text/template"

	"github.com/Masterminds/sprig/v3"
)

var (
	//go:embed dsl.tmpl
	dslTmpl string
)

type DslTmpl struct {
	PackageName string
	StructName  string
}

func (p *DslTmpl) Gen(w io.Writer) error {

	tmplExec := template.Must(template.New("h2o-dsl-tmpl").Funcs(sprig.FuncMap()).Parse(dslTmpl))

	return tmplExec.Execute(w, *p)
}

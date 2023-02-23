package types

import (
	"io"
	"text/template"

	_ "embed"

	"github.com/antlabs/h2o/pyaml"
)

//go:embed types.tmpl
var httpTypeTemplate string

func newTypeTemplate() *template.Template {
	tmpl := httpTypeTemplate
	return template.Must(template.New("h2o-http-client-type-tmpl").Parse(tmpl))
}

func Gen(t *pyaml.TypeTmpl, w io.Writer) {
	tpl := newTypeTemplate()
	tpl.Execute(w, *t)
}

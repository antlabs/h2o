package pb

import (
	"io"
	"text/template"

	_ "embed"

	"github.com/antlabs/h2o/pyaml"
)

var (
	//go:embed pb_type.tmpl
	pbTypeTemplate string
)

func newTypeTemplate() *template.Template {
	tmpl := pbTypeTemplate
	return template.Must(template.New("h2o-pb-type-tmpl").Parse(tmpl))
}

func Gen(t *pyaml.TypeTmpl, w io.Writer) {
	tpl := newTypeTemplate()
	tpl.Execute(w, *t)
}

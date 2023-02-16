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

type pbType struct {
	pyaml.TypeTmpl
	URLName string
}

func newPbType(pt pyaml.TypeTmpl, urlname string) *pbType {
	return &pbType{TypeTmpl: pt, URLName: urlname}
}

func (p *pbType) Gen(w io.Writer) {
	tpl := newTypeTemplate()
	tpl.Execute(w, *p)
}

package transport

import (
	_ "embed"
	"io"
	"text/template"

	"github.com/Masterminds/sprig/v3"
)

type goZeroMain struct {
	PackageNameSlice []string
	GoZeroBaseURL    string
	GoModName        string
}

//go:embed gozero_main.tmpl
var goZeroMainTmpl string

func newGoZeroMainTemplate() *template.Template {
	tmpl := goZeroMainTmpl
	return template.Must(template.New("h2o-transport-gozero-main-tmpl").Funcs(sprig.FuncMap()).Parse(tmpl))
}

func (g *goZeroMain) Gen(w io.Writer) error {

	tmpl := newGoZeroMainTemplate()

	return tmpl.Execute(w, *g)
}

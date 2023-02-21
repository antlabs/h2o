package transport

import (
	_ "embed"
	"io"
	"text/template"

	"github.com/Masterminds/sprig/v3"
)

type goZeroServer struct {
	GoZeroBaseURL string
	GoModName     string
	PackageName   string //包名
	Func          []Func //func
}

//go:embed gozero_server.tmpl
var goZeroServerTmpl string

func newGoZeroServerTemplate() *template.Template {
	tmpl := goZeroServerTmpl
	return template.Must(template.New("h2o-transport-gozero-server-tmpl").Funcs(sprig.FuncMap()).Parse(tmpl))
}

func (g *goZeroServer) Gen(w io.Writer) error {

	tmpl := newGoZeroServerTemplate()

	return tmpl.Execute(w, *g)
}

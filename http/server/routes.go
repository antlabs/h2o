package server

import (
	_ "embed"
	"io"
	"text/template"
)

type RoutesTmpl struct {
	AllRoute          []Routes
	GoMod             string
	AllSubPackageName []string
}

type Routes struct {
	Method         string
	Path           string
	SubPackageName string
	Handler        string
}

//go:embed routes.tmpl
var httpRoutesTemplate string

func (l *RoutesTmpl) Gen(w io.Writer) {
	tpl := func() *template.Template {
		tmpl := httpRoutesTemplate
		return template.Must(template.New("h2o-http-server-routes-tmpl").Parse(tmpl))
	}()

	tpl.Execute(w, *l)
}

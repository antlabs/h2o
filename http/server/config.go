package server

import (
	_ "embed"
	"html/template"
	"io"
)

//go:embed config.tmpl
var httpConfigTemplate string

func GenConfig(w io.Writer) {
	tpl := func() *template.Template {
		tmpl := httpConfigTemplate
		return template.Must(template.New("h2o-http-server-config-tmpl").Parse(tmpl))
	}()

	tpl.Execute(w, nil)
}

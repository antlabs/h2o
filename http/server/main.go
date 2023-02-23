package server

import (
	_ "embed"
	"html/template"
	"io"
)

type MainTmpl struct {
	GoMod string
}

//go:embed main.tmpl
var httpMainTemplate string

func Gen(m *MainTmpl, w io.Writer) {
	tpl := func() *template.Template {
		tmpl := httpMainTemplate
		return template.Must(template.New("h2o-http-server-main-tmpl").Parse(tmpl))
	}()
	tpl.Execute(w, *m)
}

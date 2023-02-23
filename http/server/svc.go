package server

import (
	_ "embed"
	"html/template"
	"io"
)

type SvcTmpl struct {
	GoMod string
}

//go:embed svc.tmpl
var SvcTemplate string

func (s *SvcTmpl) Gen(w io.Writer) {
	tpl := func() *template.Template {
		tmpl := httpMainTemplate
		return template.Must(template.New("h2o-http-server-svc-tmpl").Parse(tmpl))
	}()
	tpl.Execute(w, *s)
}

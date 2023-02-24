package server

import (
	_ "embed"
	"io"
	"text/template"
)

type HandlerTmpl struct {
	SubPackageName string
	GoMod          string
	Handler        string
	ReqName        string
	HasURL         bool
	HasHeader      bool
	HasQuery       bool
	HasJSONBody    bool
}

//go:embed handler.tmpl
var HandlerTemplate string

func (h *HandlerTmpl) Gen(w io.Writer) {
	tpl := func() *template.Template {
		tmpl := HandlerTemplate
		return template.Must(template.New("h2o-http-server-handler-tmpl").Parse(tmpl))
	}()

	tpl.Execute(w, *h)
}

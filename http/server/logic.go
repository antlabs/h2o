package server

import (
	_ "embed"
	"io"
	"text/template"
)

//go:embed logic.tmpl
var httpLogicTemplate string

type LogicTmpl struct {
	SubPackageName string
	GoMod          string
	Handler        string
	ReqName        string
	RespName       string
}

func (l *LogicTmpl) Gen(w io.Writer) {
	tpl := func() *template.Template {
		tmpl := httpLogicTemplate
		return template.Must(template.New("h2o-http-server-logic-tmpl").Parse(tmpl))
	}()

	tpl.Execute(w, *l)
}

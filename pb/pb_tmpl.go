package pb

import (
	"io"
	"text/template"
)

const (
	pbTmpl = `syntax = "proto3";

  option go_package="{{.PackageName}}";

  {{.PbType}}

  service {{.ServiceName}} {
  {{range $value := .Func}}
    rpc {{.Name}} ({{.ReqName}}) returns ({{.RespName}});
  {{end}}
  }
  `
)

type PbTmpl struct {
	ServiceName string
	PackageName string
	Func        []Func
	PbType      string //这里面都是类型定义
}

type Func struct {
	Name     string
	ReqName  string
	RespName string
}

func newFuncTemplate() *template.Template {
	tmpl := pbTmpl
	return template.Must(template.New("h2o-pb-tmpl").Parse(tmpl))
}

func (p *PbTmpl) Gen(w io.Writer) error {

	tmpl := newFuncTemplate()

	return tmpl.Execute(w, *p)
}

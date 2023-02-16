package pb

import (
	_ "embed"
	"io"
	"text/template"
)

var (
	//go:embed pb.tmpl
	pbTmpl string
)

type PbTmpl struct {
	PackageName   string //protobuf 文件中的package
	GoPackageName string //protobuf 文件中的go_papckage
	ServiceName   string //protobuf 里面的rcp ServiceName(请求) returns(响应);
	Func          []Func //
	PbType        string //这里面都是类型定义
}

type Func struct {
	Name     string //service名
	ReqName  string //service 请求参数名
	RespName string //service 响应参数名
}

func newFuncTemplate() *template.Template {
	tmpl := pbTmpl
	return template.Must(template.New("h2o-pb-tmpl").Parse(tmpl))
}

func (p *PbTmpl) Gen(w io.Writer) error {

	tmpl := newFuncTemplate()

	return tmpl.Execute(w, *p)
}

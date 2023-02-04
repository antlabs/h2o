package pb

import (
	"io"
	"text/template"

	"github.com/antlabs/h2o/pyaml"
)

const (
	pbTypeTemplate = `
  {{range $value := .ReqResp}}
// 请求包
message {{$value.Req.Name}} {
{{- if $value.Req.Query.Name}}
    {{$value.Req.Query.Name}} Query = 1; {{end}}
{{- if $value.Req.Header.Name}}
    {{$value.Req.Header.Name}} Header = 2;{{end}}
{{- if $value.Req.Body.Name }}
    {{$value.Req.Body.Name}} Body = 3; {{end}}
}

// 响应包
message {{$value.Resp.Name}} {
{{- if $value.Resp.Header.Name}}
    {{$value.Resp.Header.Name}} Header = 1;{{end}}
{{- if $value.Resp.Body.Name}}
    {{$value.Resp.Body.Name}} Body = 2;{{end}}
}

{{- if $value.Req.Query.StructType}}
// 查询字符串结构体
{{$value.Req.Query.StructType}}
{{- end}}

{{- if $value.Req.Header.StructType}}
// 请求头结构体
{{$value.Req.Header.StructType}}
{{- end}}

{{- if $value.Req.Body.StructType}}
// 请求body结构体
{{$value.Req.Body.StructType}}
{{- end}}

{{- if $value.Resp.Header.StructType}}
// 响应头结构体
{{$value.Resp.Header.StructType}}
{{- end}}

{{- if $value.Resp.Body.StructType}}
// 响应body结构体
{{$value.Resp.Body.StructType}}
{{- end}}
{{- end}}
  `
)

func newTypeTemplate() *template.Template {
	tmpl := pbTypeTemplate
	return template.Must(template.New("h2o-pb-type-tmpl").Parse(tmpl))
}

func Gen(t *pyaml.TypeTmpl, w io.Writer) {
	tpl := newTypeTemplate()
	tpl.Execute(w, *t)
}

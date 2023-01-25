package client

import (
	"io"
	"text/template"
)

// 构造函数
const (
	//htttpClientTemplate2 = `package {}`
	httpClientTemplate = `package {{.PackageName}}

import(
  "fmt"

  "github.com/guonaihong/gout"
)

type {{.StructName}} struct {
  {{- range $_, $value := .InitField}}
    {{$value}} string
  {{- end}}
}

func New() *{{.StructName}} {
  return &{{.StructName}}{}
}

// 函数
  {{- $ReceiverName := .ReceiverName}}
  {{- $StructName := .StructName}}
  {{- $url := .URL}}
{{range $_, $value := .AllFunc}}
func ({{$ReceiverName}} *{{$StructName}}) {{$value.HandlerName}}({{if $value.ReqBodyName}}req *{{$value.ReqBodyName}}{{end}}) (*{{$value.RespBodyName}}, error) {

  {{if $value.ReqBodyName}}var resp {{$value.RespBodyName}}{{end}}
  code := 0
  err := gout.{{.Method}}({{$url}}, *{{$ReceiverName}}){{if .Header}}.SetHeader(req.Header){{end}}{{if .Query}}.SetQuery(req.Query){{end}}.SetJSON(req.Body.ReqBody).BindJSON(&resp.Body).Code(&code).Do()
  if err != nil {
    return nil,err
  }

  if code != 200 {
    return nil, fmt.Errorf("{{$value.HandlerName}} code(%d) != 200", code)
  }
  return &resp, nil
}
{{end}}
 `
)

type genHTTPClient struct {
	InitField    []string //初始化的成员字段
	PackageName  string   //包名
	URL          string   //url 地址
	ReceiverName string   //接收器名
	StructName   string   //结构体
	AllFunc      []Func   //func
}

type Func struct {
	Method       string   //http方法
	Header       []string //http header
	Query        []string //htttp 查询字符串
	ReqBodyName  string   //请求body名称
	RespBodyName string   //请求body名称
	HandlerName  string   //生成的函数名
}

func newFuncTemplate() *template.Template {
	tmpl := httpClientTemplate
	return template.Must(template.New("h2o-http-client-tmpl").Parse(tmpl))
}

func (h *genHTTPClient) Gen(w io.Writer) {
	tpl := newFuncTemplate()
	tpl.Execute(w, *h)
}

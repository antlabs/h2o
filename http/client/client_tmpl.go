package client

import (
	"io"
	"text/template"

	"github.com/antlabs/h2o/model"
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
  {{- range $key, $value := .InitField}}
    {{$key}} string
  {{- end}}
}

func New() *{{.StructName}} {
  rv := &{{.StructName}}{
  {{- range $key, $value := .InitField}}
    {{$key}}:{{$value|printf "%q"}},
  {{- end}}
  }
  rv.Init()
  return rv
}

// set xxx 成员函数
{{- $ReceiverName := .ReceiverName}}
{{- $StructName := .StructName}}
{{- range $key, $value := .InitField}}
func ({{$ReceiverName}} *{{$StructName}}) Set{{$key}} ({{$key}} string) *{{$StructName}} {
  {{$ReceiverName}}.{{$key}} = {{$key}}
  return {{$ReceiverName}} 
}
{{- end}}

// 成员函数

{{- range $_, $value := .AllFunc}}
func ({{$ReceiverName}} *{{$StructName}}) {{$value.HandlerName}}({{if $value.ReqName}}req *{{$value.ReqName}}{{end}}) (*{{$value.RespName}}, error) {

  {{if $value.ReqName}}var resp {{$value.RespName}}{{end}}

  {{- range $_, $val := .DefReqHeader}}
  {{- if $val.IsString}}
  if len(req.Header.{{$val.Key}}) == 0 {
    req.Header.{{$val.Key}} = {{$val.Val|printf "%q"}}
  }
  {{- end}}
  {{- if $val.IsInt}}
  if req.Header.{{$val.Key}} == 0 {
    req.Header.{{$val.Key}} = {{$val.Val}}
  }
  {{- end}}
  {{- if $val.IsFloat64}}
  if req.Header.{{$val.Key}} == 0.0 {
    req.Header.{{$val.Key}} = {{$val.Val}}
  }
  {{- end}}
  {{- end}}

  {{- range $_, $val := .DefReqBody}}

  {{- if $val.IsString}}
  if len(req.Body.{{$val.Key}}) == 0 {
    req.Body.{{$val.Key}} = {{$val.Val|printf "%q"}}
  }
  {{- end}}
  {{- if $val.IsInt}}
  if req.Body.{{$val.Key}} == 0 {
    req.Body.{{$val.Key}} = {{$val.Val}}
  }
  {{- end}}
  {{- if $val.IsFloat64}}
  if req.Body.{{$val.Key}} == 0.0 {
    req.Body.{{$val.Key}} = {{$val.Val}}
  }
  {{- end}}
  {{- end}}

  code := 0
  err := gout.{{.Method}}({{$value.URL|printf "%q"}}, *{{$ReceiverName}}).
  RequestUse({{$ReceiverName}})
  {{- if .HaveHeader}}.
  SetHeader(req.Header)
  {{- end}}{{- if .HaveQuery}}.
  SetQuery(req.Query){{end}}{{if .HaveReqBody}}.
  {{- if .ReqWWWForm}}
  SetWWWForm(req.Body)
  {{- else}}
  SetJSON(req.Body)
  {{- end}}{{end}}.
  BindJSON(&resp.Body).
  Code(&code).
  Do()
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

	httpClientLogicTemplate = `package {{.PackageName}}

  import(
    "net/http"
  )

  func ({{.ReceiverName}} *{{.StructName}}) Init() {
  }

  func ({{.ReceiverName}} *{{.StructName}}) ModifyRequest(req *http.Request) (err error) {
    return
  }
`
)

type ClientTmpl struct {
	InitField    map[string]string //初始化的成员字段
	PackageName  string            //包名
	ReceiverName string            //接收器名
	StructName   string            //结构体
	AllFunc      []Func            //func
}

type Func struct {
	URL          string //url 地址
	Method       string //http方法 GET POST DELETE之类的
	DefReqHeader []model.KeyVal[string, string]
	DefReqBody   []model.KeyVal[string, string]
	HaveHeader   bool   //有http header
	HaveQuery    bool   //有查询字符串
	HaveReqBody  bool   //有请求body
	ReqWWWForm   bool   //www form 编码
	ReqName      string //函数请求参数名
	RespName     string //函数响应参数名
	HandlerName  string //生成的函数名
}

func newFuncTemplate() *template.Template {
	tmpl := httpClientTemplate
	return template.Must(template.New("h2o-http-client-tmpl").Parse(tmpl))
}

func newFuncLogicTemplate() *template.Template {
	tmpl := httpClientLogicTemplate
	return template.Must(template.New("h2o-http-client-logic-tmpl").Parse(tmpl))
}

// 生成业务逻辑
func (h *ClientTmpl) GenLogic(w io.Writer) {
	tpl := newFuncLogicTemplate()
	tpl.Execute(w, *h)
}

func (h *ClientTmpl) Gen(w io.Writer) {
	tpl := newFuncTemplate()
	tpl.Execute(w, *h)
}

package client

import (
	"io"
	"text/template"
)

const (
	httpTypeTemplate = `
  {{range $value := .ReqResp}}
    type {{$value.Req.Name}} struct {
      Query {{$value.Req.QueryName}}
      Body {{$value.Req.BodyName}}
      Header {{$value.Req.HeaderName}}
    }

    type {{$value.Resp.Name}} struct {
      Header {{$value.Resp.Header.Name}}
      Body {{$value.Resp.Body.Name}}
    }

    // 查询字符串结构体
    {{$value.Req.Query.StructType}}
    // 请求头结构体
    {{$value.Req.Header.StructType}}
    // 请求body结构体
    {{$value.Req.Body.StructType}}

    // 响应头结构体
    {{$value.Resp.Header.StructType}}
    // 响应body结构体
    {{$value.Resp.Body.StructType}}
  {{end}}
  `
)

type Query struct {
	Name       string
	StructType []byte
}

type Body struct {
	Name       string
	StructType []byte
}

type Header struct {
	Name       string
	StructType []byte
}

type Req struct {
	Name   string //请求的结构体名
	Query  Query  //Query string 名
	Body   Body   //body名
	Header Header //header名
}

type Resp struct {
	Name   string
	Body   Body   //响应body结构体名
	Header Header //响应header结构体名
}

type ReqResp struct {
	Req  Req  //请求
	Resp Resp //响应
}
type TypeTmpl struct {
	ReqResp []ReqResp
}

func newTypeTemplate() *template.Template {
	tmpl := httpTypeTemplate
	return template.Must(template.New("h2o-http-client-type-tmpl").Parse(tmpl))
}

func (t *TypeTmpl) Gen(w io.Writer) {
	tpl := newTypeTemplate()
	tpl.Execute(w, *t)
}

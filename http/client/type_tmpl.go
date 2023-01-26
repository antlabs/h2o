package client

import (
	"io"
	"text/template"
)

const (
	httpTypeTemplate = `
package {{.PackageName}}

  {{range $value := .ReqResp}}
    type {{$value.Req.Name}} struct {
      {{if $value.Req.Query.Name}} Query {{$value.Req.QueryName}} {{end}}
      {{if $value.Req.Body.Name }} Body {{$value.Req.Body.Name}} {{end}}
      {{if $value.Req.Header.Name}} Header {{$value.Req.Header.Name}} {{end}}
    }

    type {{$value.Resp.Name}} struct {
      {{if $value.Resp.Header.Name}} Header {{$value.Resp.Header.Name}} {{end}}
      {{if $value.Resp.Body.Name}} Body {{$value.Resp.Body.Name}} {{end}}
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
	StructType string
}

type Body struct {
	Name       string
	StructType string
}

type Header struct {
	Name       string
	StructType string
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
	PackageName string
	ReqResp     []ReqResp
}

func newTypeTemplate() *template.Template {
	tmpl := httpTypeTemplate
	return template.Must(template.New("h2o-http-client-type-tmpl").Parse(tmpl))
}

func (t *TypeTmpl) Gen(w io.Writer) {
	tpl := newTypeTemplate()
	tpl.Execute(w, *t)
}

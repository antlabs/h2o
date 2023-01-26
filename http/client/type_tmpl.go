package client

import (
	"io"
	"text/template"
)

const (
	httpTypeTemplate = `
  type {{.Req.Name}} struct {
    Query {{.Req.QueryName}}
    Body {{.Req.BodyName}}
    Header {{.Req.HeaderName}}
  }

  type {{.Resp.Name}} struct {
    Header {{.Resp.HeaderName}}
    Body {{.Resp.BodyName}}
  }

  // 查询字符串结构体
  {{.ReqQueryStruct}}
  // 请求头结构体
  {{.ReqHeaderStruct}}
  // 请求body结构体
  {{.ReqBodyStruct}}

  // 响应头结构体
  {{.RespHeaderStruct}}
  // 响应body结构体
  {{.RespBodyStruct}}
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
	Body   Body   //响应body结构体名
	Header Header //响应header结构体名
}

type TypeTmpl struct {
	Req  Req  //请求
	Resp Resp //响应
}

func newTypeTemplate() *template.Template {
	tmpl := httpClientTemplate
	return template.Must(template.New("h2o-http-client-type-tmpl").Parse(tmpl))
}

func (t *TypeTmpl) Gen(w io.Writer) {
	tpl := newFuncTemplate()
	tpl.Execute(w, *t)
}

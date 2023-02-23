package client

import (
	_ "embed"
	"io"
	"text/template"

	"github.com/antlabs/h2o/model"
)

// 构造函数
var (

	//go:embed client.tmpl
	httpClientTemplate string

	//go:embed logic.tmpl
	httpClientLogicTemplate string
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
	URLTemplate  bool   //url 启用模板语法
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

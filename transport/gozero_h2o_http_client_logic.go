package transport

import (
	_ "embed"
	"io"
	"text/template"
	//"github.com/Masterminds/sprig"
)

var (
	//go:embed gozero_h2o_http_client_logic.tmpl
	goZeroH2oHTTPClientLogic string
)

type transportGoZeroHTTPClientTmpl struct {
	SvcName           string
	PackageName       string //包名
	URLStruct         string //请求体名
	GoZeroBaseURL     string //
	HTTPClientBaseURL string //
	Func              Func   //func
}

type Func struct {
	RpcName  string //service名
	ReqName  string //service 请求参数名
	RespName string //service 响应参数名
}

func newFuncTemplate() *template.Template {
	tmpl := goZeroH2oHTTPClientLogic
	return template.Must(template.New("h2o-transport-gozero-http-client-tmpl").Parse(tmpl))
}

func (g *transportGoZeroHTTPClientTmpl) Gen(w io.Writer) error {

	tmpl := newFuncTemplate()

	return tmpl.Execute(w, *g)
}

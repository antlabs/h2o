package codemsg

import (
	"html/template"
	"io"

	_ "embed"
)

//go:embed grpc_codemsg.tmpl
var grpcCodemsgTmpl string

// 自定义
type GrpcCodeMsgTmpl struct {
	PkgName      string // 包名
	TypeName     string
	MsgName      string // 修改Message 字段的名字
	Args         string // os.Args[2:]
	OriginalName string //
	AllVariable  []Value
}

func newGrpcFuncTemplate() *template.Template {
	tmpl := grpcCodemsgTmpl
	return template.Must(template.New("h2o-codemsg-grpc-tmpl").Parse(tmpl))
}

func (c *GrpcCodeMsgTmpl) Gen(w io.Writer) error {
	tpl := newGrpcFuncTemplate()

	return tpl.Execute(w, *c)
}

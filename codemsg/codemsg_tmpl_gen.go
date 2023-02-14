package codemsg

import (
	_ "embed"
	"html/template"
	"io"

	"github.com/gobeam/stringy"
)

//go:embed codemsg.tmpl
var codemsgTmpl string

type CodeMsgTmpl struct {
	PkgName           string // 包名
	TypeName          string
	CodeMsgStructName string // CodeMsg{Code int, Message string} 结构体的名字
	CodeName          string // 修改Code字段的名字
	MsgName           string // 修改Message 字段的名字
	MsgTagName        string
	CodeTagName       string
	Args              string // os.Args[2:]
	OriginalName      string //
	AllVariable       []Value
}

func newFuncTemplate() *template.Template {
	tmpl := codemsgTmpl
	return template.Must(template.New("h2o-codemsg-tmpl").Parse(tmpl))
}

func (c *CodeMsgTmpl) Gen(w io.Writer) error {
	tpl := newFuncTemplate()

	c.CodeTagName = stringy.New(c.CodeName).SnakeCase("?", "").ToLower()
	c.MsgTagName = stringy.New(c.MsgName).SnakeCase("?", "").ToLower()
	return tpl.Execute(w, *c)
}

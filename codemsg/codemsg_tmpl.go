package codemsg

import (
	"html/template"
	"io"

	"github.com/gobeam/stringy"
)

const codemsgTmpl = `

	// Code generated by "h2o codemsg {{.Args}}"; DO NOT EDIT."
	
	package {{.PkgName}}
	
  import (
    "strings"
    "encoding/json"
  )
	
	type {{.CodeMsgStructName}} struct {
    {{.CodeName}}    {{.TypeName}} "json:\"{{.CodeTagName}}\"" 
    {{.MsgName}} string  "json:\"{{.MsgTagName}}\""
	}

	func (x *{{.CodeMsgStructName}}) Error() string {
    all, _ := json.Marshal(x) 
    var b strings.Builder 
    b.Write(all) 
    return b.String()
	}

	func New{{.CodeMsgStructName}}(code {{.TypeName}}) error {
    return &CodeMsg{
      {{.CodeName}}: code,
      {{.MsgName}}: code.String(),
	  }
	}

  {{- $CodeMsgStructName := .CodeMsgStructName}}
  var (
  {{range $_, $value := .AllVariable}}
		ErrCodeMsg{{$value.OriginalName}} error = New{{$CodeMsgStructName}}({{.OriginalName}}) //{{$value.Name}}
  {{end}}
  )
  `

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

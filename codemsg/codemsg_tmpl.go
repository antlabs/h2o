package codemsg

import (
	"html/template"
	"io"
)

const codemsgTmpl = `

	// Code generated by "h2o codemsg {{.Args}}"; DO NOT EDIT."
	
	package {{.PkgName}}
	
  import (
    "strings"
    "encoding/json"
  )
	
	type {{.CodeMsgName}} struct {
    {{.CodeName}}    {{.TypeName}} "json:\"{{.CodeName}}\"" 
    {{.MessageName}} string  "json:\"{{.MessageName}}\""
	}

	func (x *{{.CodeMsgName}}) Error() string {
    all, _ := json.Marshal(x) 
    var b strings.Builder 
    b.Write(all) 
    return b.String()
	}

	func New{{.CodeMsgName}}(code {{.TypeName}}) error {
    return &CodeMsg{
      {{.CodeName}}: code,
      {{.MessageName}}: code.String(),
	  }
	}

  {{- $CodeMsgName := .CodeMsgName}}
  var (
  {{range $_, $value := .AllVariable}}
		ErrCodeMsg{{$value.OriginalName}} error = New{{$CodeMsgName}}({{.OriginalName}}) //{{$value.Name}}
  {{end}}
  )
  `

type CodeMsgTmpl struct {
	PkgName      string // 包名
	TypeName     string
	CodeMsgName  string   // CodeMsg{Code int, Message string} 结构体的名字
	CodeName     string   // 修改Code字段的名字
	MessageName  string   // 修改Message 字段的名字
	Args         []string // os.Args[2:]
	OriginalName string   //
	AllVariable  []Value
}

func newFuncTemplate() *template.Template {
	tmpl := codemsgTmpl
	return template.Must(template.New("h2o-codemsg-tmpl").Parse(tmpl))
}

func (c *CodeMsgTmpl) Gen(w io.Writer) {
	tpl := newFuncTemplate()
	tpl.Execute(w, *c)
}

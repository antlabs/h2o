package curl

import (
	_ "embed"
	"io"
	"text/template"
)

//go:embed curl_short.tmpl
var shortTmpl string

//go:embed curl_long.tmpl
var longTmpl string

type curlTmpl struct {
	Method   string
	Header   []string
	URL      string
	Data     string
	FormData []string
	Long     bool //不参与生成
}

func (c *curlTmpl) Gen(w io.Writer) {

	tmpl := func() *template.Template {
		tmpl := shortTmpl
		if c.Long {
			tmpl = longTmpl
		}

		return template.Must(template.New("generate-curl").Parse(tmpl))
	}()

	tmpl.Execute(w, *c)
}

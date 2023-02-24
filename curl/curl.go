package curl

import (
	"bytes"
	"encoding/json"
	"fmt"

	h2ourl "github.com/antlabs/h2o/internal/url"
	"github.com/antlabs/h2o/parser"
)

type Curl struct {
	File []string `clop:"short;long;greedy" usage:"Parsing dst files" valid:"required"`
	Long bool     `clop:"short;long" usage:"short or long option curl command"`
}

func (c *Curl) SubMain() {

	for _, f := range c.File {

		dsl, err := parser.Parser(f)
		if err != nil {
			fmt.Printf("h2o.HTTP.parser %s\n", err)
			return
		}

		for _, m := range dsl.Multi {

			var all string
			if m.Req.Body != nil {

				allBytes, err := json.Marshal(m.Req.Body)
				if err != nil {
					panic(err.Error())
				}
				all = string(allBytes)
			}

			var out bytes.Buffer

			urlStr := m.Req.URL
			if m.Req.Template.URL {
				urlStr = h2ourl.TakeURL(m.Req.URL)
			}

			tmpl := curlTmpl{Method: m.Req.Method, Header: m.Req.Header, URL: urlStr, Long: c.Long, Data: all}
			tmpl.Gen(&out)
			fmt.Printf("%s\n\n", out.String())
		}
	}

}

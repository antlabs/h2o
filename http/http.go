package http

import (
	"fmt"
	"os"
	"strings"

	"github.com/antlabs/h2o/http/client"
	"github.com/antlabs/h2o/parser"
	"github.com/antlabs/tostruct/option"
	"github.com/antlabs/tostruct/url"
)

// 生成客户端代码
// h2o http -f ./testdata/dst.yaml -g client

// TODO
// 生成服务端代码
// h2o http -f ./testdata/dst.yaml -g server

// TODO
// 生成客户端和服务端代码
// h2o http -f ./testdata/dst.yaml -g client server
type HTTP struct {
	File []string `clop:"short;long;greedy" usage:"Parsing dst files" valid:"required"`
	Gen  []string `clop:"short;long;greedy" usage:"Generate client or server code" valid:"required"`
}

func (h *HTTP) SubMain() {
	for _, f := range h.File {

		c, err := parser.Parser(f)
		if err != nil {
			fmt.Printf("http.parser %s\n", err)
			return
		}

		tmpl := client.ClientTmpl{PackageName: c.Package,
			InitField:    c.Init.Resp.Field,
			StructName:   c.Init.Resp.Name,
			ReceiverName: string(c.Init.Resp.Name[0]),
		}

		for _, h := range c.Muilt {

			handler := h.Handler
			if pos := strings.Index(handler, "."); pos != -1 {
				handler = handler[pos+1:]
			}

			// http header
			var header []string
			if len(h.Req.Header) > 0 {
				header = make([]string, 0, len(h.Req.Header)*2)
				for _, v := range h.Req.Header {
					pos := strings.Index(v, ":")
					if pos == -1 {
						continue
					}
					header = append(header, v[:pos])
					header = append(header, v[pos+1:])
				}
			}

			queryName := h.Req.Name + "Query"
			all, err := url.Marshal(h.URL, option.WithStructName(queryName), option.WithTagName("query"))
			tmpl.AllFunc = append(tmpl.AllFunc, client.Func{
				HandlerName: handler,
				Method:      h.Req.Method,
				URL:         h.URL,
				ReqName:     h.Req.Name,
				RespName:    h.Resp.Name,
				HaveHeader:  len(h.Req.Header) > 0,
			})

			_ = err
			_ = all
		}

		tmpl.Gen(os.Stdout)
	}
}

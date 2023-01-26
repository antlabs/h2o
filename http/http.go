package http

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/antlabs/h2o/http/client"
	"github.com/antlabs/h2o/parser"
	"github.com/antlabs/tostruct/header"
	"github.com/antlabs/tostruct/json"
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

func getBody(name string, bodyData any) (body client.Body, err error) {

	body.Name = name + "Body"

	switch v := bodyData.(type) {
	case map[string]any:
		body.StructType, err = json.Marshal(v, option.WithStructName(body.Name), option.WithTagName("json"))
	case []any:
		body.StructType, err = json.Marshal(v, option.WithStructName(body.Name), option.WithTagName("json"))
	}

	return
}

func getHeader(name string, headerArray []string) (htmpl client.Header, err error) {

	// http header
	if len(headerArray) == 0 {
		return
	}

	htmpl.Name = name
	var hmap http.Header
	for _, v := range headerArray {
		pos := strings.Index(v, ":")
		if pos == -1 {
			continue
		}
		hmap.Set(v[:pos], v[pos+1:])
	}

	htmpl.StructType, err = header.Marshal(hmap, option.WithStructName(htmpl.Name), option.WithTagName("header"))
	if err != nil {
		return
	}

	return
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

		tmplType := client.TypeTmpl{}

		for _, h := range c.Muilt {

			handler := h.Handler
			if pos := strings.Index(handler, "."); pos != -1 {
				handler = handler[pos+1:]
			}

			var query client.Query
			if len(h.Req.URL) > 0 {
				query.Name = h.Req.Name + "Query"
				all, err := url.Marshal(h.Req.URL, option.WithStructName(query.Name), option.WithTagName("query"))
				query.StructType = all
				if err != nil {
					fmt.Printf("marshal query string fail:%s\n", err)
					return
				}
			}

			reqHeader, err := getHeader(h.Req.Name+"Header", h.Req.Header)
			if err != nil {
				fmt.Printf("get request header:%s\n", err)
				return
			}

			respHeader, err := getHeader(h.Resp.Name+"Header", h.Resp.Header)
			if err != nil {
				fmt.Printf("get request body:%s\n", err)
				return
			}

			reqBody, err := getBody(h.Req.Name, h.Req.Body)
			if err != nil {
				fmt.Printf("get request body:%s\n", err)
				return
			}

			respBody, err := getBody(h.Resp.Name, h.Resp.Body)
			if err != nil {
				fmt.Printf("get response body:%s\n", err)
				return
			}

			tmplType.ReqResp = append(tmplType.ReqResp, client.ReqResp{
				Req: client.Req{
					Name:   h.Req.Name,
					Query:  query,
					Header: reqHeader,
					Body:   reqBody,
				},
				Resp: client.Resp{
					Name:   h.Resp.Name,
					Header: respHeader,
					Body:   respBody,
				},
			})

			// TODO 检查handler名
			tmpl.AllFunc = append(tmpl.AllFunc, client.Func{
				HandlerName: handler,
				Method:      h.Req.Method,
				URL:         h.Req.URL,
				ReqName:     h.Req.Name,
				RespName:    h.Resp.Name,
				HaveHeader:  len(h.Req.Header) > 0,
			})
		}

		tmpl.Gen(os.Stdout)
		tmplType.Gen(os.Stdout)
	}
}

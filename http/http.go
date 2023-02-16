package http

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/antlabs/h2o/http/client"
	"github.com/antlabs/h2o/internal/save"
	"github.com/antlabs/h2o/model"
	"github.com/antlabs/h2o/parser"
	"github.com/antlabs/h2o/pyaml"
	"github.com/antlabs/tostruct/option"
	"github.com/antlabs/tostruct/url"
)

// 1.生成客户端代码, OK
// h2o http -f ./testdata/dst.yaml -g client

// TODO
// 2.生成服务端代码
// h2o http -f ./testdata/dst.yaml -g server

// 3.TODO
// 生成客户端和服务端代码
// h2o http -f ./testdata/dst.yaml -g client server
type HTTP struct {
	File []string `clop:"short;long;greedy" usage:"Parsing dst files" valid:"required"`
	Gen  []string `clop:"short;long;greedy" usage:"Generate client or server code" valid:"required"`
	Dir  string   `clop:"short;long" usage:"gen dir" default:"."`
}

// 入口函数
func (h *HTTP) SubMain() {
	for _, f := range h.File {

		c, err := parser.Parser(f)
		if err != nil {
			fmt.Printf("h2o.HTTP.parser %s\n", err)
			return
		}

		tmplClient := client.ClientTmpl{PackageName: c.Package,
			InitField:    c.Init.Resp.Field,
			StructName:   c.Init.Resp.Name,
			ReceiverName: strings.ToLower(string(c.Init.Resp.Name[0])),
		}

		tmplType := pyaml.TypeTmpl{PackageName: c.Package}

		for _, h := range c.Muilt {

			h.ModifyHandler()
			handler := h.Handler
			// 去除前面的包名, TODO 没有如果点就是普通函数
			if handler == "" {
				panic("handler name is empty")
			}

			h.Req.Name, h.Resp.Name = h.GetReqName(), h.GetRespName()

			if h.Req.Name == "" {
				panic("req name is empty")
			}
			var query pyaml.Query
			if len(h.Req.URL) > 0 {
				query.Name = h.Req.Name + "Query"
				if pos := strings.Index(h.Req.URL, "?"); pos != -1 {
					// dst.yaml里面的字符串带{}两种花括号。直接使用解析这样的url会报错。给个默认正确的host，就可以得到query string
					urlStr := "www.qq.com?" + h.Req.URL[pos+1:]
					all, err := url.Marshal(urlStr, option.WithStructName(query.Name), option.WithTagName("query"))
					h.Req.URL = h.Req.URL[:pos] //删除原始url里面的query string, 这里已经通过SetQuery里面传过来了
					query.StructType = string(all)
					if err != nil {
						fmt.Printf("marshal query string fail:%s\n", err)
						return
					}
				} else {
					query.Name = ""
				}
			}

			reqHeader, defReqHeader, respHeader, _, err := pyaml.GetHeader(h)
			if err != nil {
				return
			}

			reqBody, defReqBody, respBody, err := pyaml.GetBody(h, false)

			tmplType.ReqResp = append(tmplType.ReqResp, pyaml.ReqResp{
				Req: pyaml.Req{
					Name:   h.Req.Name,
					Query:  query,
					Header: reqHeader,
					Body:   reqBody,
				},
				Resp: pyaml.Resp{
					Name:   h.Resp.Name,
					Header: respHeader,
					Body:   respBody,
				},
			})

			// TODO 检查handler名, 不能有重复
			tmplClient.AllFunc = append(tmplClient.AllFunc, client.Func{
				HandlerName:  handler,
				Method:       h.Req.Method,
				URL:          h.Req.URL,
				URLTemplate:  h.Req.Template.URL,
				ReqName:      h.Req.Name,
				RespName:     h.Resp.Name,
				DefReqHeader: defReqHeader,
				DefReqBody:   defReqBody,
				ReqWWWForm:   h.Req.Encode.Body == model.WWWForm,
				HaveQuery:    len(query.StructType) > 0,
				HaveHeader:   len(h.Req.Header) > 0,
				HaveReqBody:  h.Req.Body != nil,
			})
		}

		dir := save.Mkdir(h.Dir, tmplClient.PackageName)
		save.TmplFile(getFuncName(dir, tmplClient.PackageName), true, func() []byte {
			var buf bytes.Buffer
			tmplClient.Gen(&buf)
			return buf.Bytes()
		})

		save.TmplFile(getTypeName(dir, tmplClient.PackageName), true, func() []byte {
			var typeBuf bytes.Buffer
			client.Gen(&tmplType, &typeBuf)
			return typeBuf.Bytes()
		})
		save.TmplFile(getLogicName(dir, tmplClient.PackageName), true, func() []byte {
			var buf bytes.Buffer
			tmplClient.GenLogic(&buf)
			return buf.Bytes()
		})

	}
}

func getFuncName(dir string, packageName string) string {
	return dir + "/" + packageName + ".go"
}

func getTypeName(dir string, packageName string) string {
	return dir + "/" + packageName + "_type.go"
}

func getLogicName(dir string, packageName string) string {
	return dir + "/" + packageName + "_logic.go"
}

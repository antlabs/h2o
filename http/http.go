package http

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/antlabs/deepcopy"
	"github.com/antlabs/h2o/http/client"
	"github.com/antlabs/h2o/http/server"
	"github.com/antlabs/h2o/http/types"
	"github.com/antlabs/h2o/internal/gomod"
	"github.com/antlabs/h2o/internal/save"
	h2ourl "github.com/antlabs/h2o/internal/url"
	"github.com/antlabs/h2o/model"
	"github.com/antlabs/h2o/parser"
	"github.com/antlabs/h2o/pyaml"
	"github.com/antlabs/pcurl"
	"github.com/antlabs/pcurl/body"
	"github.com/antlabs/tostruct/option"
	"github.com/antlabs/tostruct/url"
)

// 错误处理风格
// h2o 是代码生成器，有错误可以直接panic，这和服务端处理逻辑不一样

// 1.生成客户端代码, OK
// h2o http -f ./testdata/dst.yaml --client

// 2.生成服务端代码
// h2o http -f ./testdata/dst.yaml --server

// 3.生成客户端和服务端代码
// h2o http -f ./testdata/dst.yaml --client --server
type HTTP struct {
	File   []string `clop:"short;long;greedy" usage:"Parsing dst files" valid:"required"`
	Client bool     `clop:"short;long" usage:"gen http client code"`
	Server bool     `clop:"short;long" usage:"gen http server code"`
	Dir    string   `clop:"short;long" usage:"gen dir" default:"."`
	Debug  bool     `clop:"long" usage:"debug mode"`
}

// 入口函数
func (h *HTTP) SubMain() {
	goModName := gomod.GetGoModuleName(h.Dir)
	routes := server.RoutesTmpl{GoMod: goModName}
	for _, f := range h.File {

		c, err := parser.Parser(f)
		if err != nil {
			fmt.Printf("h2o.HTTP.parser %s, filename(%s)\n", err, f)
			return
		}

		tmplClient := client.ClientTmpl{
			PackageName:  c.Package,
			InitField:    c.Init.RvStruct.Field,
			StructName:   c.Init.RvStruct.Name,
			ReceiverName: strings.ToLower(string(c.Init.RvStruct.Name[0])),
		}

		tmplClientType := pyaml.TypeTmpl{PackageName: c.Package}
		routes.AllSubPackageName = append(routes.AllSubPackageName, c.Package)

		hp := h
		logicDir := ""
		handlerDir := ""
		if h.Server {
			logicDir = save.MkdirAndClean(getLogicPrefix(hp.Dir, c.Package))
			handlerDir = save.MkdirAndClean(getHandlerPrefix(hp.Dir, c.Package))
		}

		for _, h := range c.Multi {

			if len(h.Req.Curl) > 0 {
				reqObj, err := pcurl.ParseAndObj(h.Req.Curl)
				if err != nil {
					panic(err.Error()) // 提前报错，让用户修复下数据
				}

				if err := deepcopy.Copy(&h.Req, reqObj).Do(); err != nil {
					panic(err.Error())
				}
			}

			if h.Resp.Body != nil {
				if s, ok := h.Resp.Body.(string); ok {
					_, o, err := body.Unmarshal([]byte(s))
					if err != nil {
						panic(err.Error())
					}
					h.Resp.Body = o
				}
			}
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
			if hp.Debug {
				fmt.Printf("resp:%#v\n", h.Resp.Body)
				fmt.Printf("req.Query:%#v\n", h.Req.URL)
			}

			if len(h.Req.URL) > 0 {
				query.Name = h.Req.Name + "Query"
				if pos := strings.Index(h.Req.URL, "?"); pos != -1 {
					// dst.yaml里面的字符串带{}两种花括号。直接使用解析这样的url会报错。给个默认正确的host，就可以得到query string
					urlStr := "www.qq.com?" + h.Req.URL[pos+1:]
					all, err := url.Marshal(urlStr, option.WithStructName(query.Name), option.WithTagName("query"))
					h.Req.URL = h.Req.URL[:pos] // 删除原始url里面的query string, 这里已经通过SetQuery里面传过来了
					query.StructType = string(all)
					if err != nil {
						fmt.Printf("marshal query string fail:%s\n", err)
						return
					}
				} else {
					query.Name = ""
				}
			}

			hasQuery := len(query.StructType) > 0
			hasHeader := len(h.Req.Header) > 0
			hasJSONBody := h.Req.Body != nil

			if hp.Server {

				save.TmplFile(getLogicName(logicDir, handler), true, func() []byte {
					var buf bytes.Buffer
					logicTmpl := server.LogicTmpl{SubPackageName: c.Package, GoMod: goModName, Handler: handler, ReqName: h.Req.Name, RespName: h.Resp.Name}
					logicTmpl.Gen(&buf)
					return buf.Bytes()
				})

				save.TmplFile(getHandlerName(handlerDir, handler), true, func() []byte {
					var buf bytes.Buffer
					handlerTmpl := server.HandlerTmpl{
						SubPackageName: c.Package,
						GoMod:          goModName,
						Handler:        handler,
						ReqName:        h.Req.Name,
						HasQuery:       hasQuery,
						HasHeader:      hasHeader,
						HasJSONBody:    hasJSONBody,
					}
					handlerTmpl.Gen(&buf)
					return buf.Bytes()
				})

			}

			reqHeader, defReqHeader, respHeader, _, err := pyaml.GetHeader(h)
			if err != nil {
				return
			}

			reqBody, defReqBody, respBody, err := pyaml.GetBody(h, false)

			tmplClientType.ReqResp = append(tmplClientType.ReqResp, pyaml.ReqResp{
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
				HaveQuery:    hasQuery,
				HaveHeader:   hasHeader,
				HaveReqBody:  hasJSONBody,
				HaveRespBody: len(respBody.StructType) > 0,
			})

			if hp.Server {
				routes.AllRoute = append(routes.AllRoute,
					server.Routes{
						Method:         h.Req.Method,
						Path:           h2ourl.TakePath(h.Req.URL, h.Req.Template.URL),
						SubPackageName: c.Package,
						Handler:        handler,
					})
			}
		}

		if h.Client {
			dir := save.Mkdir(h.Dir, tmplClient.PackageName)
			save.TmplFile(getClientFuncName(dir, tmplClient.PackageName), true, func() []byte {
				var buf bytes.Buffer
				tmplClient.Gen(&buf)
				return buf.Bytes()
			})
			// 保存客户端结构体定义
			save.TmplFile(getClientTypeName(dir, tmplClient.PackageName), true, func() []byte {
				var typeBuf bytes.Buffer
				types.Gen(&tmplClientType, &typeBuf)
				return typeBuf.Bytes()
			})
			save.TmplFile(getClientLogicName(dir, tmplClient.PackageName), true, func() []byte {
				var buf bytes.Buffer
				tmplClient.GenLogic(&buf)
				return buf.Bytes()
			})
		}

		if h.Server {
			// 保存服务端结构体定义
			dir := save.MkdirAndClean(getServerTypePrefix(h.Dir, tmplClient.PackageName))
			save.TmplFile(getServerTypeName(dir, tmplClient.PackageName), true, func() []byte {
				var typeBuf bytes.Buffer
				types.Gen(&tmplClientType, &typeBuf)
				return typeBuf.Bytes()
			})
		}
	}

	if h.Server {
		// 保存main.go 服务入口文件
		goModLastName := gomod.GetGoModuleLastName(h.Dir)
		dir := save.MkdirAndClean(getServerPrefixMain(h.Dir, goModLastName))

		// 保存至 xx.go并且格式化
		save.TmplFile(getServerMainName(dir, goModLastName), true, func() []byte {
			var typeBuf bytes.Buffer
			server.Gen(&server.MainTmpl{GoMod: goModName, GoModLastName: goModLastName}, &typeBuf)
			return typeBuf.Bytes()
		})

		// svc, 保存至servercontext.go并且格式化
		dir = save.MkdirAndClean(getServerSvcPrefix(h.Dir))
		save.TmplFile(getServerSvcName(dir), true, func() []byte {
			var typeBuf bytes.Buffer
			(&server.SvcTmpl{GoMod: goModName}).Gen(&typeBuf)
			return typeBuf.Bytes()
		})

		// config， 保存到config.go并且格式化
		dir = save.MkdirAndClean(getServerConfigPrefix(h.Dir))
		save.TmplFile(getServerConfigName(dir), true, func() []byte {
			var typeBuf bytes.Buffer
			server.GenConfig(&typeBuf)
			return typeBuf.Bytes()
		})

		// routes 保存至的routes.go并且格式化
		save.TmplFile(getRoutesName(getRoutesPrefix(h.Dir)), true, func() []byte {
			var typeBuf bytes.Buffer
			routes.Gen(&typeBuf)
			return typeBuf.Bytes()
		})

		// routes 保存至的routes.go并且格式化
		dir = save.MkdirAndClean(getEtcPrefix(h.Dir))
		save.TmplFile(getEtcName(dir, goModLastName), false, func() []byte {
			var typeBuf bytes.Buffer
			g := server.ConfigYAML{GoModLastName: goModLastName}
			g.Gen(&typeBuf)
			return typeBuf.Bytes()
		})
	}
}

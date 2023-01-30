package http

import (
	"bytes"
	stdjson "encoding/json"
	"fmt"
	"go/format"
	"net/http"
	"os"
	"path"
	"sort"
	"strings"

	"github.com/antlabs/h2o/http/client"
	"github.com/antlabs/h2o/model"
	"github.com/antlabs/h2o/parser"
	"github.com/antlabs/tostruct/header"
	"github.com/antlabs/tostruct/json"
	"github.com/antlabs/tostruct/name"
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
	Dir  string   `clop:"short;long" usage:"gen dir" default:"."`
}

func getBody(bodyName string, bodyData any, newType map[string]string, encode model.Encode, bodyDefKey []string) (
	body client.Body,
	rvDefaultBody []model.KeyVal[string, string],
	err error) {

	body.Name = bodyName + "Body"

	tagName := "json"
	if encode.Body == model.WWWForm {
		tagName = "form"
	}

	getVal := make(map[string]any)

	for _, v := range bodyDefKey {
		getVal[v] = ""
	}

	var data []byte
	switch v := bodyData.(type) {
	case map[string]any:
		data, err = json.Marshal(v, option.WithStructName(body.Name), option.WithTagName(tagName), option.WithSpecifyType(newType),
			option.WithGetRawValue(getVal))
	case []any:
		data, err = json.Marshal(v, option.WithStructName(body.Name), option.WithTagName(tagName), option.WithSpecifyType(newType),
			option.WithGetRawValue(getVal))
	default:
		body.Name = ""
	}

	if len(getVal) > 0 {
		rvDefaultBody = make([]model.KeyVal[string, string], 0, len(getVal))
		for k, v := range getVal {
			fieldName, _ := name.GetFieldAndTagName(k)
			rvDefaultBody = append(rvDefaultBody, model.KeyVal[string, string]{Key: fieldName, Val: fmt.Sprint(v), RawVal: v})

		}
		sort.Slice(rvDefaultBody, func(i, j int) bool {
			return rvDefaultBody[i].Key < rvDefaultBody[i].Key
		})
	}

	body.StructType = string(data)
	return
}

func getHeader(headerName string, headerSlice []string, defaultHeader []string) (
	htmpl client.Header,
	rvDefaultHeader []model.KeyVal[string, string],
	err error) {

	// http header
	if len(headerSlice) == 0 {
		return
	}

	hmap := sliceToHTTPHeader(headerSlice)
	htmpl.Name = headerName

	getVal := make(map[string]any)

	for _, v := range defaultHeader {
		_, ok := hmap[v]
		if !ok {
			continue
		}

		getVal[v] = ""
		//fieldName, _ := name.GetFieldAndTagName(v)
		//rvDefaultHeader = append(rvDefaultHeader, model.KeyVal[string, string]{Key: fieldName, Val: hv[0]})
	}

	var data []byte
	data, err = header.Marshal(hmap, option.WithStructName(htmpl.Name),
		option.WithTagName("header"), option.WithTagNameFromKey(), option.WithGetRawValue(getVal))
	if err != nil {
		return
	}

	if len(getVal) > 0 {
		rvDefaultHeader = make([]model.KeyVal[string, string], 0, len(getVal))
		for k, v := range getVal {
			fieldName, _ := name.GetFieldAndTagName(k)
			rvDefaultHeader = append(rvDefaultHeader, model.KeyVal[string, string]{Key: fieldName, Val: fmt.Sprint(v), RawVal: v})
		}
		sort.Slice(rvDefaultHeader, func(i, j int) bool {
			return rvDefaultHeader[i].Key < rvDefaultHeader[i].Key
		})
	}

	htmpl.StructType = string(data)
	return
}

func exists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}

	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

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
			ReceiverName: string(c.Init.Resp.Name[0]),
		}

		tmplType := client.TypeTmpl{PackageName: c.Package}

		for _, h := range c.Muilt {

			handler := h.Handler
			// 去除前面的包名, TODO 没有如果点就是普通函数
			if pos := strings.Index(handler, "."); pos != -1 {
				handler = handler[pos+1:]
			}

			if handler == "" {
				panic("handler name is empty")
			}

			if h.Req.Name == "" {
				h.Req.Name = handler + "Req"
			}

			if h.Resp.Name == "" {
				h.Resp.Name = handler + "Resp"
			}

			var query client.Query
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

			reqHeader, defReqHeader, err := getHeader(h.Req.Name+"Header", h.Req.Header, h.Req.UseDefault.Header)
			if err != nil {
				fmt.Printf("get request header:%s, %v\n", err, h.Req.Header)
				return
			}

			respHeader, _, err := getHeader(h.Resp.Name+"Header", h.Resp.Header, nil)
			if err != nil {
				fmt.Printf("get request body:%s\n", err)
				return
			}

			reqBody, defReqBody, err := getBody(h.Req.Name, h.Req.Body, h.Req.NewType, h.Req.Encode, h.Req.UseDefault.Body)
			if err != nil {
				fmt.Printf("get request body:%s\n", err)
				return
			}

			respBody, _, err := getBody(h.Resp.Name, h.Resp.Body, h.Resp.NewType, model.Encode{}, nil)
			if err != nil {
				fmt.Printf("get response body:%s \n", err)
				all, _ := stdjson.Marshal(h.Resp.Body)
				fmt.Println(string(all), err, h.Resp.Body == nil)
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

			// TODO 检查handler名, 不能有重复
			tmplClient.AllFunc = append(tmplClient.AllFunc, client.Func{
				HandlerName:  handler,
				Method:       h.Req.Method,
				URL:          h.Req.URL,
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

		dir := h.Dir + "/" + tmplClient.PackageName
		dir = path.Clean(dir)
		os.MkdirAll(dir, 0755)

		saveTempFile(getFuncName(dir, tmplClient.PackageName), func() []byte {
			var buf bytes.Buffer
			tmplClient.Gen(&buf)
			return buf.Bytes()
		})

		saveTempFile(getTypeName(dir, tmplClient.PackageName), func() []byte {
			var typeBuf bytes.Buffer
			tmplType.Gen(&typeBuf)
			return typeBuf.Bytes()
		})
		saveTempFile(getLogicName(dir, tmplClient.PackageName), func() []byte {
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

func saveTempFile(fileName string, getTmpl func() []byte) {

	if b, _ := exists(fileName); !b {

		buf := getTmpl()
		fmtType, err := format.Source(buf)
		if err != nil {
			fmt.Printf("%s fail:%s\n", fileName, err)
			os.Stdout.Write(buf)
			return
		}

		//os.Stdout.Write(fmtType)
		os.WriteFile(fileName, fmtType, 0644)
	} else {
		fmt.Printf("%s 已经存在，忽略\n", fileName)
	}
}

func sliceToHTTPHeader(headerSlice []string) http.Header {

	hmap := make(http.Header)
	for _, v := range headerSlice {
		pos := strings.Index(v, ":")
		if pos == -1 {
			continue
		}

		val := v[pos+1:]
		if len(val) == 0 {
			continue
		}
		if val[0] == ' ' {
			val = val[1:]
		}
		hmap.Set(v[:pos], val)
	}
	return hmap
}

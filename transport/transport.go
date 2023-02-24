package transport

import (
	"bytes"
	"fmt"
	"os"

	"github.com/antlabs/h2o/internal/gomod"
	"github.com/antlabs/h2o/internal/save"
	"github.com/antlabs/h2o/parser"
)

const (
	goZero     = "go-zero"
	httpClient = "h2o-http-client"
)

// 这是主模块
// TransportTmpl模块，就是选用一种服务端接入数据，然后选用一种客户端把数据发走
type Transport struct {
	// 来自哪个服务端
	FromType string `clop:"long" usage:"server template type" default:"go-zero"`
	// 推到哪个客户端
	ToType string `clop:"long" usage:"client template type" default:"h2o-http-client"`
	// 设置服务端包的基本url
	FromBaseURL string `clop:"--from-base-url" usage:"from base url"`
	// 设置客户端包的基本url
	ToBaseURL string `clop:"--to-base-url" usage:"to base url"`
	// 解析
	File []string `clop:"short;long;greedy" usage:"Parsing dst files" valid:"required"`
	// 生成目录
	Dir string `clop:"short;long" usage:"gen dir" default:"."`
}

// 入口
func (t *Transport) SubMain() {

	gomod := gomod.GetGoModuleName(t.Dir)
	tmplMain := goZeroMain{GoModName: gomod, GoZeroBaseURL: t.FromBaseURL}
	tmplServer := goZeroServer{GoZeroBaseURL: t.FromBaseURL, GoModName: gomod}
	for _, f := range t.File {

		//var out bytes.Buffer
		c, err := parser.Parser(f)
		if err != nil {
			fmt.Printf("h2o.HTTP.parser %s\n", err)
			return
		}

		tmplMain.PackageNameSlice = append(tmplMain.PackageNameSlice, c.Package)
		tmplServer.PackageName = c.Package
		for _, h := range c.Multi {

			h.ModifyHandler()

			//protobuf.Marshal()
			h.Req.Name, h.Resp.Name = h.GetReqName(), h.GetRespName()

			if h.Handler == "" {
				panic("h.Handler is empty")
			}

			if t.FromType == goZero && t.ToType == httpClient {

				tmpl := transportGoZeroHTTPClientTmpl{
					SvcName:           gomod + "/internal/svc",
					URLStruct:         c.Init.Resp.Name,
					PackageName:       c.Package,
					GoZeroBaseURL:     t.FromBaseURL,
					HTTPClientBaseURL: t.ToBaseURL}

				tmpl.Func = Func{
					RpcName:  h.Handler,
					ReqName:  h.Req.Name,
					RespName: h.Resp.Name}

				tmplServer.Func = append(tmplServer.Func, tmpl.Func)

				save.TmplFile(getLogicName(t.Dir, c.Package, h.Handler), true, func() []byte {
					var typeBuf bytes.Buffer
					tmpl.Gen(&typeBuf)
					return typeBuf.Bytes()
				})
			}
		}
		save.TmplFile(getServerName(t.Dir, tmplServer.PackageName), true, func() []byte {
			var typeBuf bytes.Buffer
			tmplServer.Gen(&typeBuf)
			return typeBuf.Bytes()
		})
		tmplServer.Func = nil
	}

	// 生成main.go
	save.TmplFile(getMainName(t.Dir, gomod), true, func() []byte {
		var typeBuf bytes.Buffer
		tmplMain.Gen(&typeBuf)
		return typeBuf.Bytes()
	})
}

func getServerName(dir string, packageName string) string {
	return fmt.Sprintf("%s/internal/server/%[2]s/%[2]sserver.go", dir, packageName)

}
func getMainName(dir string, gomod string) string {
	return dir + "/" + gomod + ".go"
}

func getLogicName(dir string, packageName string, handler string) string {
	prefix := fmt.Sprintf("%s/internal/logic/%s", dir, packageName)

	os.MkdirAll(prefix, 0755)
	return fmt.Sprintf("%s/%slogic.go", prefix, handler)
}

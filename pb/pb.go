package pb

import (
	"bytes"
	"fmt"

	"github.com/antlabs/h2o/internal/save"
	"github.com/antlabs/h2o/parser"
	"github.com/antlabs/h2o/pyaml"
	"github.com/antlabs/tostruct/option"
	//"github.com/antlabs/tostruct/protobuf"
)

type Pb struct {
	File []string `clop:"short;long;greedy" usage:"Parsing dst files" valid:"required"`
	Dir  string   `clop:"short;long" usage:"gen dir" default:"."`
}

func (b *Pb) SubMain() {

	for _, f := range b.File {

		var out bytes.Buffer
		c, err := parser.Parser(f)
		if err != nil {
			fmt.Printf("h2o.HTTP.parser %s\n", err)
			return
		}

		tmplType := pyaml.TypeTmpl{PackageName: c.Package}
		packageName := c.Package
		if c.Protobuf.GoPackage != "" {
			packageName = c.Protobuf.GoPackage
		}
		tmplPb := PbTmpl{PackageName: packageName, ServiceName: c.Package}
		for _, h := range c.Muilt {

			h.ModifyHandler()

			//protobuf.Marshal()
			h.Req.Name, h.Resp.Name = h.GetReqName(), h.GetRespName()

			if h.Handler == "" {
				panic("h.Handler is empty")
			}

			reqHeader, _, respHeader, _, err := pyaml.GetHeader(h, option.WithProtobuf())
			if err != nil {
				return
			}

			reqBody, _, respBody, err := pyaml.GetBody(h, true)

			tmplType.ReqResp = append(tmplType.ReqResp, pyaml.ReqResp{
				Req: pyaml.Req{
					Name:   h.Req.Name,
					Header: reqHeader,
					Body:   reqBody,
				},
				Resp: pyaml.Resp{
					Name:   h.Resp.Name,
					Header: respHeader,
					Body:   respBody,
				},
			})

			var typeOut bytes.Buffer
			Gen(&tmplType, &typeOut)
			tmplPb.PbType = typeOut.String()
			tmplPb.Func = append(tmplPb.Func, Func{Name: h.Handler, ReqName: h.Req.Name, RespName: h.Resp.Name})

		}

		dir := save.Mkdir(b.Dir, c.Package)
		//fmt.Println("dir###", dir, b.Dir, c.Package)
		save.TmplFile(getProtoBuf(dir, c.Package), false, func() []byte {
			tmplPb.Gen(&out)
			return out.Bytes()
			//os.Stdout.Write(out.Bytes())
		})
	}
}

func getProtoBuf(name string, pacakge string) string {
	return fmt.Sprintf("%s/%s.proto", name, pacakge)
}

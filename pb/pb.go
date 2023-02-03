package pb

import (
	"bytes"
	"fmt"
	"os"

	"github.com/antlabs/h2o/parser"
	"github.com/antlabs/h2o/pyaml"
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
		tmplPb := PbTmpl{PackageName: c.Package, ServiceName: c.Package}
		for _, h := range c.Muilt {

			h.ModifyHandler()

			//protobuf.Marshal()
			reqName, respName := h.GetReqName(), h.GetRespName()

			if h.Handler == "" {
				panic("h.Handler is empty")
			}

			reqHeader, _, respHeader, _, err := pyaml.GetHeader(h)
			if err != nil {
				return
			}

			reqBody, _, respBody, err := pyaml.GetBody(h)

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
			tmplPb.Func = append(tmplPb.Func, Func{Name: h.Handler, ReqName: reqName, RespName: respName})
		}

		tmplPb.Gen(&out)
		os.Stdout.Write(out.Bytes())
	}
}

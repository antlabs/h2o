package client

import (
	"fmt"
	"os"
	"testing"
)

type testData struct {
}

func (t *testData) Write(p []byte) (n int, err error) {
	return
}

func Test_newFuncTmpl(t *testing.T) {

	g := &ClientTmpl{
		InitField: map[string]any{"Host": "default.host",
			"Org_name": "default.org_name",
			"App_name": "default.App_name",
			"Onlykey":  ""},
		URL:          fmt.Sprintf("%q", `https://{{.Host}}/{{.Org_name}}/{{.App_name}}/users`),
		PackageName:  "users",
		ReceiverName: "u",
		StructName:   "Users",
		AllFunc: []Func{{
			Method:       "POST",
			Header:       []string{"h1", "h1value", "h2", "h2value"},
			HandlerName:  "CreateUser",
			ReqBodyName:  "UserReq",
			RespBodyName: "UserResp",
		}}}
	g.Gen(os.Stderr)
	//g.Gen(os.Stdout)
}

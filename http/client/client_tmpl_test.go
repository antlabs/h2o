package client

import (
	"os"
	"testing"

	"github.com/antlabs/h2o/model"
)

type testData struct {
}

func (t *testData) Write(p []byte) (n int, err error) {
	return
}

func Test_newFuncTmpl(t *testing.T) {

	g := &ClientTmpl{
		InitField: map[string]string{"Host": "default.host",
			"Org_name": "default.org_name",
			"App_name": "default.App_name",
			"Onlykey":  ""},
		PackageName:  "users",
		ReceiverName: "u",
		StructName:   "Users",
		AllFunc: []Func{{
			DefReqHeader: []model.KeyVal[string, string]{{Key: "Accept", Val: "json"}},
			URL:          `https://{{.Host}}/{{.Org_name}}/{{.App_name}}/users`,
			Method:       "POST",
			HaveHeader:   true,
			HaveQuery:    true,
			HaveReqBody:  true,
			HandlerName:  "CreateUser",
			ReqName:      "UserReq",
			RespName:     "UserResp",
		}}}
	g.Gen(os.Stderr)
	//g.Gen(os.Stdout)
}

func Test_newFuncTmpl2(t *testing.T) {

	g := &ClientTmpl{
		InitField: map[string]string{"Host": "default.host",
			"Org_name": "default.org_name",
			"App_name": "default.App_name",
			"Onlykey":  ""},
		PackageName:  "users",
		ReceiverName: "u",
		StructName:   "Users",
		AllFunc: []Func{{
			//DefReqHeader: []model.KeyVal[string, string]{{Key: "Accept", Val: "json"}},
			URL:         `https://{{.Host}}/{{.Org_name}}/{{.App_name}}/users`,
			Method:      "POST",
			HaveHeader:  true,
			HaveQuery:   true,
			HandlerName: "CreateUser",
			ReqName:     "UserReq",
			RespName:    "UserResp",
		}}}
	g.Gen(os.Stderr)
	//g.Gen(os.Stdout)
}

func Test_newFuncTmpl_wwwform(t *testing.T) {

	g := &ClientTmpl{
		InitField: map[string]string{"Host": "default.host",
			"Org_name": "default.org_name",
			"App_name": "default.App_name",
			"Onlykey":  ""},
		PackageName:  "users",
		ReceiverName: "u",
		StructName:   "Users",
		AllFunc: []Func{{
			ReqWWWForm: true,
			//DefReqHeader: []model.KeyVal[string, string]{{Key: "Accept", Val: "json"}},
			URL:         `https://{{.Host}}/{{.Org_name}}/{{.App_name}}/users`,
			Method:      "POST",
			HaveHeader:  true,
			HaveQuery:   true,
			HaveReqBody: true,
			HandlerName: "CreateUser",
			ReqName:     "UserReq",
			RespName:    "UserResp",
		}}}
	g.Gen(os.Stderr)
	//g.Gen(os.Stdout)
}

func Test_newFuncTmpl_BodyDef(t *testing.T) {

	g := &ClientTmpl{
		InitField: map[string]string{"Host": "default.host",
			"Org_name": "default.org_name",
			"App_name": "default.App_name",
			"Onlykey":  ""},
		PackageName:  "users",
		ReceiverName: "u",
		StructName:   "Users",
		AllFunc: []Func{{
			ReqWWWForm: true,
			DefReqBody: []model.KeyVal[string, string]{
				{
					Key: ".grant_type",
					Val: "client_credentials",
				},
			},
			//DefReqHeader: []model.KeyVal[string, string]{{Key: "Accept", Val: "json"}},
			URL:         `https://{{.Host}}/{{.Org_name}}/{{.App_name}}/users`,
			Method:      "POST",
			HaveHeader:  true,
			HaveQuery:   true,
			HaveReqBody: true,
			HandlerName: "CreateUser",
			ReqName:     "UserReq",
			RespName:    "UserResp",
		}}}
	g.Gen(os.Stderr)
	//g.Gen(os.Stdout)
}

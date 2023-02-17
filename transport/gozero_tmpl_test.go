package transport

import (
	"os"
	"testing"
)

func Test_GoZeroHTTPClient(t *testing.T) {
	tmpl := transportGoZeroHTTPClientTmpl{PackageName: "users", GoZeroBaseURL: "gitlab.xx/server", HTTPClientBaseURL: "gitlab.xx/client"}
	tmpl.Func = Func{RpcName: "CreateUser",
		ReqName:  "CreateUserReq",
		RespName: "CreateUserResp"}
	tmpl.Gen(os.Stdout)
}

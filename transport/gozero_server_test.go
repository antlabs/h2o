package transport

import (
	"bytes"
	"go/format"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGoZero_Server(t *testing.T) {
	var out bytes.Buffer
	gs := goZeroServer{GoZeroBaseURL: "gitlab/base",
		GoModName:   "admin",
		PackageName: "hello",
		Func: []Func{{RpcName: "Create",
			ReqName:  "CreateReq",
			RespName: "CreateResp"}}}
	gs.Gen(&out)

	_, err := format.Source(out.Bytes())
	assert.NoError(t, err, out.String())
	//os.Stdout.Write(all)
}

package server

import (
	"os"
	"testing"
)

func TestLogic(t *testing.T) {
	l := LogicTmpl{SubPackageName: "upload", GoMod: "main", Handler: "Upload", ReqName: "UploadReq", RespName: "UploadResp"}
	l.Gen(os.Stdout)
}

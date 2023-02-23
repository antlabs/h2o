package server

import (
	"os"
	"testing"
)

func TestHandler(t *testing.T) {
	handler := HandlerTmpl{
		SubPackageName: "users",
		GoMod:          "main",
		Handler:        "CreateUser",
		ReqName:        "CreateUserReq",
		HasURL:         true,
		HasHeader:      true,
		HasQuery:       true,
		HasJSONBody:    true,
	}

	handler.Gen(os.Stdout)
}

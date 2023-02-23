package server

import (
	"os"
	"testing"
)

func TestRoutes(t *testing.T) {
	r := RoutesTmpl{
		AllRoute: []Routes{
			{Method: "GET", Path: "/im/hello", SubPackageName: "users", Handler: "GetHello"},
			{Method: "POST", Path: "/im/hello", SubPackageName: "users", Handler: "CreateHello"},
		},
	}
	r.Gen(os.Stdout)
}

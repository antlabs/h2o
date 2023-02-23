package server

import (
	"os"
	"testing"
)

func TestMain(m *testing.T) {
	Gen(&MainTmpl{GoMod: "xx"}, os.Stdout)
}

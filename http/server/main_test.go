package server

import (
	"os"
	"testing"
)

func TestMain(m *testing.T) {
	Gen(&MainTmpl{GoMod: "githug.com/g/xx", GoModLastName: "xx"}, os.Stdout)
}

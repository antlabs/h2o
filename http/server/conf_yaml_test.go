package server

import (
	"os"
	"testing"
)

func Test_ConfYAML(t *testing.T) {

	g := ConfigYAML{GoModLastName: "hello"}
	g.Gen(os.Stdout)
}

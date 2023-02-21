package transport

import (
	"bytes"
	"go/format"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_GoZeroMain(t *testing.T) {
	g := goZeroMain{GoZeroBaseURL: "gitlab.com/xx",
		PackageNameSlice: []string{"users", "metadata", "apptoken", "contacts", "usertoken"},
		GoModName:        "admin"}
	var out bytes.Buffer

	g.Gen(&out)

	_, err := format.Source(out.Bytes())
	assert.NoError(t, err, out.String())
}

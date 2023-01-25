package pyaml

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Parse(t *testing.T) {

	p, err := Parse("../testdata/dst.yaml")
	assert.NoError(t, err)

	for _, d := range p.Muilt {
		fmt.Printf("%s, %T, %s, %T\n", d.Req.Name, d.Req.Body, d.Resp.Name, d.Resp.Body)
	}
}

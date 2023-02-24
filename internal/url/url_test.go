package url

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type testTakeURLCase struct {
	Data       string
	Need       string
	IsTemplate bool
}

func TestTakeURL(t *testing.T) {
	for _, v := range []testTakeURLCase{
		{Data: "https://host/im/api", Need: "/im/api", IsTemplate: false},
		{Data: "https://{{.Host}}/im/{{.Action}}", Need: "/im/:Action", IsTemplate: true},
	} {
		got := TakePath(v.Data, v.IsTemplate)
		assert.Equal(t, got, v.Need)
	}
}

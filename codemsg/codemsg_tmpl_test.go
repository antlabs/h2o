package codemsg

import (
	"bytes"
	"fmt"
	"go/format"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_CodeMsg(t *testing.T) {

	c := CodeMsgTmpl{
		PkgName:           "demo",
		CodeMsgStructName: "CodeMsg", // CodeMsg{Code int, Message string} 结构体的名字
		CodeName:          "Code",    // 修改Code字段的名字
		MsgName:           "Message", // 修改Message 字段的名字
		TypeName:          "ErrNo",
		Args:              "--code-msg --linecomment    --type    ErrNo    ./err.go", // os.Args[2:]
		OriginalName:      "demo",                                                    //
		StringMethod:      "String",
		AllVariable:       []Value{{OriginalName: "aa", Name: "bb"}},
	}

	var buf bytes.Buffer
	c.Gen(&buf)
	fmt.Println(buf.String())
	out, err := format.Source(buf.Bytes())
	assert.NoError(t, err)
	fmt.Println(string(out))
}

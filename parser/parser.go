package parser

import (
	"errors"
	"strings"

	"github.com/antlabs/h2o/model"
	"github.com/antlabs/h2o/pyaml"
)

// 不支持的文件类型
var ErrNotSupport = errors.New("not support file type")

func Parser(fileName string) (c *model.Config, err error) {
	if strings.HasSuffix(fileName, ".yaml") {
		return pyaml.Parse(fileName)
	}
	return nil, ErrNotSupport
}

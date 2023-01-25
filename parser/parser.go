package parser

import (
	"strings"

	"github.com/antlabs/h2o/model"
	"github.com/antlabs/h2o/pyaml"
)

func Parser(fileName string) (c *model.Config, err error) {
	if strings.HasSuffix(fileName, ".yaml") {
		return pyaml.Parse(fileName)
	}
	return
}

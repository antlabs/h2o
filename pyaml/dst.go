package pyaml

import (
	"os"

	"github.com/antlabs/h2o/model"
	"gopkg.in/yaml.v3"
)

func Parse(fileName string) (c *model.Config, err error) {
	all, err := os.ReadFile(fileName)
	if err != nil {
		return nil, err
	}
	c = &model.Config{}
	err = yaml.Unmarshal(all, &c)
	if err != nil {
		return nil, err
	}
	return
}

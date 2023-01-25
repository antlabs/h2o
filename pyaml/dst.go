package pyaml

import (
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Muilt   []Muilt `yaml:"muilt"`
	Package string  `yaml:"package"`
}

type Muilt struct {
	Handler  string `yaml:"handler"`
	Req      Req    `yaml:"req"`
	ReqName  string `yaml:"reqName"`
	Resp     Resp   `yaml:"resp"`
	RespName string `yaml:"respName"`
	URL      string `yaml:"url"`
	RespType any    `yaml:"respType"`
}

type Req struct {
	Body   any      `yaml:"body"`
	Header []string `yaml:"header"`
	Method string   `yaml:"method"`
}

type Resp struct {
	Body any `yaml:"body"`
}

func Parse(fileName string) (c *Config, err error) {
	all, err := os.ReadFile(fileName)
	if err != nil {
		return nil, err
	}
	c = &Config{}
	err = yaml.Unmarshal(all, &c)
	if err != nil {
		return nil, err
	}
	return
}

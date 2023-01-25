package model

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

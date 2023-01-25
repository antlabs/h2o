package model

type Config struct {
	Init    Init    `yaml:"init"`
	Muilt   []Muilt `yaml:"muilt"`
	Package string  `yaml:"package"`
}

type Init struct {
	Handler string   `yaml:"handler"`
	Req     InitReq  `yaml:"req"`
	Resp    InitResp `yaml:"resp"`
}

type InitReq struct {
	Field map[string]string `yaml:"field"`
}

type InitResp struct {
	Name  string            `yaml:"name"`
	Field map[string]string `yaml:"field"`
}

type Muilt struct {
	Handler string `yaml:"handler"`
	Req     Req    `yaml:"req"`
	Resp    Resp   `yaml:"resp"`
	URL     string `yaml:"url"`
}

type Req struct {
	Name   string   `yaml:"name"`
	Body   any      `yaml:"body"`
	Header []string `yaml:"header"`
	Method string   `yaml:"method"`
}

type Resp struct {
	Name string `yaml:"name"`
	Body any    `yaml:"body"`
}

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
	Name string            `yaml:"name"`
	Type map[string]string `yaml:"type"`
}

type InitResp struct {
	Name string            `yaml:"name"`
	Type map[string]string `yaml:"type"`
}

type Muilt struct {
	Handler  string            `yaml:"handler"`
	Req      Req               `yaml:"req"`
	ReqName  string            `yaml:"reqName"`
	Resp     Resp              `yaml:"resp"`
	RespName string            `yaml:"respName"`
	URL      string            `yaml:"url"`
	RespType map[string]string `yaml:"respType"`
}

type Req struct {
	Body   any      `yaml:"body"`
	Header []string `yaml:"header"`
	Method string   `yaml:"method"`
}

type Resp struct {
	Body any `yaml:"body"`
}

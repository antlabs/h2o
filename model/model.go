package model

import "strings"

type Config struct {
	Init     Init     `yaml:"init"`
	Multi    []Multi  `yaml:"multi"`
	Package  string   `yaml:"package"`
	Protobuf Protobuf `yaml:"protobuf"`
}

type Protobuf struct {
	Package   string `yaml:"package"`
	GoPackage string `yaml:"go_package"`
}

type Init struct {
	Handler  string   `yaml:"handler"`
	Args     InitReq  `yaml:"args"`
	RvStruct InitResp `yaml:"rvStruct"`
}

type InitReq struct {
	Field map[string]string `yaml:"field"`
}

type InitResp struct {
	Name  string            `yaml:"name"`
	Field map[string]string `yaml:"field"`
}

type Multi struct {
	Handler    string `yaml:"handler"`
	structName string
	Req        Req  `yaml:"req"`
	Resp       Resp `yaml:"resp"`
}

func (m *Multi) ModifyHandler() {

	if pos := strings.Index(m.Handler, "."); pos != -1 {
		m.structName = m.Handler[:pos]
		m.Handler = m.Handler[pos+1:]
	}
}

func (m *Multi) GetReqName() string {
	if m.Req.Name == "" {
		return m.Handler + "Req"
	}
	return m.Req.Name
}

func (m *Multi) GetRespName() string {
	if m.Resp.Name == "" {
		return m.Handler + "Resp"
	}
	return m.Resp.Name
}

type UseDefault struct {
	Header []string `yaml:"header"`
	Body   []string `yaml:"body"`
}

type Encode struct {
	Body string `yaml:"body"`
}

type Template struct {
	URL bool `yaml:"url"`
}

type UsePtr struct {
	Body []string `yaml:"body"`
}

type Req struct {
	Encode          Encode            `yaml:"encode"`
	URL             string            `yaml:"url"`
	Curl            string            `yaml:"curl"`
	UsePtr          UsePtr            `yaml:"usePtr"`
	Template        Template          `yaml:"template"`
	Name            string            `yaml:"name"`
	NewType         map[string]string `yaml:"newType"`
	NewProtobufType map[string]string `yaml:"newProtobufType"`
	Body            any               `yaml:"body"`
	Header          []string          `yaml:"header"`
	Method          string            `yaml:"method"`
	UseDefault      UseDefault        `yaml:"useDefault"`
}

type Resp struct {
	Name            string            `yaml:"name"`
	NewType         map[string]string `yaml:"newType"`
	NewProtobufType map[string]string `yaml:"newProtobufType"`
	Body            any               `yaml:"body"`
	Header          []string          `yaml:"header"`
}

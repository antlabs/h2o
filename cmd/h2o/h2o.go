package main

import (
	"fmt"
	"os"

	"github.com/antlabs/h2o/ast"
	"github.com/antlabs/h2o/codemsg"
	"github.com/antlabs/h2o/curl"
	"github.com/antlabs/h2o/http"
	"github.com/antlabs/h2o/internal/file"
	"github.com/antlabs/h2o/pb"
	"github.com/antlabs/h2o/tmpl"
	"github.com/antlabs/h2o/transport"
	"github.com/antlabs/tostruct/json"
	"github.com/antlabs/tostruct/option"
	"github.com/antlabs/tostruct/protobuf"
	"github.com/antlabs/tostruct/yaml"
	"github.com/guonaihong/clop"
)

// 从json生成结构体
type JSON struct {
	NotInline  bool   `clop:"short;long" usage:"controls whether the generated structure is split"`
	From       string `clop:"short;long" usage:"which file to open from? If it is -, it means reading from stdin"`
	StructName string `clop:"short;long" usage:"Control the name of the generated structure"`
	yaml       bool
	protobuf   bool
}

func (j *JSON) SubMain() {
	var opt []option.OptionFunc
	if j.NotInline {
		opt = append(opt, option.WithNotInline())
	}

	if len(j.StructName) > 0 {
		opt = append(opt, option.WithStructName(j.StructName))
	}

	var all []byte
	var err error

	all, err = file.ReadFile(j.From)
	if err != nil {
		fmt.Printf("open %s fail:%s\n", j.From, err)
		return
	}

	if j.yaml {
		all, err = yaml.Marshal(all, opt...)
	} else if j.protobuf {
		all, err = protobuf.Marshal(all, opt...)
	} else {
		all, err = json.Marshal(all, opt...)
	}

	if err != nil {
		fmt.Printf("write stdout fail:%s\n", err)
	}

	os.Stdout.Write(all)

}

// 从yaml生成结构体
type YAML struct {
	JSON
}

func (y *YAML) SubMain() {
	y.JSON.yaml = true
	y.JSON.SubMain()

}

type Protobuf struct {
	JSON
}

func (p *Protobuf) SubMain() {
	p.JSON.protobuf = true
	p.JSON.SubMain()
}

// 主命令
type H2O struct {
	// 子命令，入口函数是SubMain
	JSONStruct JSON `clop:"subcommand" usage:"Generate structure from json"`
	// 子命令，入口函数是SubMain
	YAMLStruct YAML `clop:"subcommand" usage:"Generate structure from yaml"`
	// 子命令，入口函数是SubMain
	ProtobufMsg Protobuf `clop:"subcommand" usage:"Generate protobuf message from json or yaml"`
	//子命令， 入口是SubMain
	CodeMsg codemsg.CodeMsg `clop:"subcommand" usage:"Generate code in codemsg format from constants"`
	// 子命令，生成http客户端代码和http服务端代码(TODO)
	HTTP http.HTTP `clop:"subcommand" usage:"gen http code(client and server)"`
	//子命令，生成protobuf
	PB pb.Pb `clop:"subcommand" usage:"gen protobuf code"`
	//transport
	Transport transport.Transport `clop:"subcommand" usage:"transport"`
	// curl， 生成curl命令
	Curl curl.Curl `clop:"subcommand" usage:"gen curl command"`
	// ast , 打印golang 的ast
	Ast ast.Ast `clop:"subcommand" usage:"print ast"`
	// tmpl, 初始化各种模板
	Tmpl tmpl.Tmpl `clop:"subcommand" usage:"tmplate init"`
}

func main() {
	var h H2O
	clop.Bind(&h)
}

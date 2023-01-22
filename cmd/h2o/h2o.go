package main

import (
	"fmt"
	"io"
	"os"

	"github.com/antlabs/tostruct/json"
	"github.com/antlabs/tostruct/option"
	"github.com/antlabs/tostruct/yaml"
	"github.com/guonaihong/clop"
)

type JSON struct {
	NotInline  bool   `clop:"short;long" usage:"controls whether the generated structure is split"`
	From       string `clop:"short;long" usage:"which file to open from? If it is -, it means reading from stdin"`
	StructName string `clop:"short;long" usage:"Control the name of the generated structure"`
	yaml       bool
}

func (j *JSON) SubMain() {
	var opt []option.OptionFunc
	if j.NotInline {
		opt = append(opt, option.WithNotInline())
	}

	var all []byte
	var err error
	if j.From == "-" {
		all, err = io.ReadAll(os.Stdin)
	} else {
		all, err = os.ReadFile(j.From)
	}
	if err != nil {
		fmt.Printf("open %s fail:%s\n", j.From, err)
		return
	}

	if j.yaml {
		all, err = yaml.Marshal(all, opt...)
	} else {
		all, err = json.Marshal(all, opt...)
	}

	if err != nil {
		fmt.Printf("write stdout fail:%s\n", err)
	}

	os.Stdout.Write(all)

}

type YAML struct {
	JSON
}

func (y *YAML) SubMain() {
	y.JSON.yaml = true
	y.JSON.SubMain()

}

type H2O struct {
	JSON JSON `clop:"subcommand" usage:"Generate structure from json"`
	YAML YAML `clop:"subcommand" usage:"Generate structure from yaml"`
}

func main() {
	var h H2O
	clop.Bind(&h)
}

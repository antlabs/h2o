package tmpl

import (
	"bytes"
	"io"
	"os"
	"path"
	"strings"
)

// tmpl子命令主要初始化，一些模板

type Tmpl struct {
	PackageName string `clop:"short;long" usage:"package name" default:"demo"`
	SaveFile    bool   `clop:"short;long" usage:"save to file"`
	Dir         string `clop:"short;long" usage:"dir" default:"."`
	Dsl         bool   `clop:"long" usage:"init dsl tmpl"`
}

func (t *Tmpl) initDsl() {
	os.MkdirAll(t.Dir, 0755)
	d := DslTmpl{PackageName: t.PackageName, StructName: strings.Title(t.PackageName)}
	var out bytes.Buffer
	d.Gen(&out)
	if t.SaveFile {
		os.WriteFile(path.Clean(t.Dir+t.PackageName), out.Bytes(), 0644)
		return
	}

	io.Copy(os.Stdout, &out)
}

func (t *Tmpl) SubMain() {
	if t.Dsl {
		t.initDsl()
	}
}

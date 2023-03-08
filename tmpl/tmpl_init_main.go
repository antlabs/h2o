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
	PackageName string `clop:"short;long" usage:"package name"`
	SaveFile    string `clop:"short;long" usage:"save to file"`
	Dsl         bool   `clop:"long" usage:"init dsl tmpl"`
}

func (t *Tmpl) initDsl() {
	pos := strings.LastIndex(t.SaveFile, "/")
	dir := ""
	if pos != -1 {
		dir = t.SaveFile[:pos]
		if dir == "" {
			dir = "./"
		}
		packageName := t.SaveFile[pos+1:]
		if strings.HasSuffix(packageName, ".yaml") {
			packageName = packageName[:len(packageName)-5]
		}
		if t.PackageName == "" {
			t.PackageName = packageName
		}
	}
	if t.PackageName == "" {
		t.PackageName = "demo"
	}

	d := DslTmpl{PackageName: t.PackageName, StructName: strings.Title(t.PackageName)}

	var out bytes.Buffer
	d.Gen(&out)
	if len(t.SaveFile) > 0 {

		os.MkdirAll(dir, 0755)
		os.WriteFile(path.Clean(t.SaveFile), out.Bytes(), 0644)
		return
	}

	io.Copy(os.Stdout, &out)
}

func (t *Tmpl) SubMain() {
	if t.Dsl {
		t.initDsl()
	}
}

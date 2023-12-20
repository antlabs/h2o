package codemsg

import (
	"bytes"
	"fmt"
	"go/format"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

const (
	takeToMapStart = "@TakeCodeToMap("
	takeToMapEnd   = ")"
)

// text样子是 @take_to_map(map名字)
func parseTakeCodeToMap(text string) (newText, mapName string, err error) {
	start := strings.Index(text, takeToMapStart)
	if start == -1 {
		return text, "", nil
	}

	end := strings.Index(text[start+len(takeToMapStart):], takeToMapEnd)
	if end == -1 {
		return text, "", nil
	}

	mapName = text[start+len(takeToMapStart) : start+len(takeToMapStart)+end]

	old := text[start : start+len(takeToMapStart)+end+1]
	newText = strings.ReplaceAll(text, old, "")
	return newText, mapName, nil
}

func saveTakeFileMap(dir, packageName, types string, mapName string, m map[int]bool) {
	baseName := fmt.Sprintf("%s_take_to_map_%s.go", types, mapName)
	fileName := filepath.Join(dir, strings.ToLower(baseName))

	const tmpl = `package {{.PackageName}}
		var {{.MapName}} = {{printf "%#v" .MyMap}}`
	tmplParsed, err := template.New("example").Parse(tmpl)
	if err != nil {
		panic(err)
	}

	// 渲染模板到一个缓冲区
	var buf bytes.Buffer

	data := struct {
		PackageName string
		MapName     string
		MyMap       map[int]bool
	}{
		MapName:     mapName,
		MyMap:       m,
		PackageName: packageName,
	}
	err = tmplParsed.Execute(&buf, data)
	if err != nil {
		panic(err)
	}
	sourceCode, err := format.Source(buf.Bytes())
	if err != nil {
		panic(string(buf.Bytes()))
	}

	err = os.WriteFile(fileName, sourceCode, 0o644)
	if err != nil {
		panic(err)
	}
}

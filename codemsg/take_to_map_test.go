package codemsg

import (
	"bytes"
	"fmt"
	"testing"
	"text/template"
)

type test_TakeCodeToMap struct {
	data     string
	needData string
	mapName  string
}

func Test_TakeCodeToMap(t *testing.T) {
	t.Run("test 成功的情况中", func(t *testing.T) {
		for _, v := range []test_TakeCodeToMap{
			{"@TakeCodeToMap(mapName)", "", "mapName"},
			{"@TakeCodeToMap(mapName) 鉴权失败", " 鉴权失败", "mapName"},
			{"鉴权失败 @TakeCodeToMap(mapName)", "鉴权失败 ", "mapName"},
			{"鉴权失败 @TakeCodeToMap(mapName) 鉴权失败的情况有3种", "鉴权失败  鉴权失败的情况有3种", "mapName"},
		} {
			newData, newMapName, err := parseTakeCodeToMap(v.data)
			if err != nil {
				t.Errorf("parseTakeCodeToMap(%s) error(%v)", v.data, err)
			}

			if newData != v.needData {
				t.Errorf("parseTakeCodeToMap[%s] newData[%s] needData[%s]", v.data, newData, v.needData)
			}

			if newMapName != v.mapName {
				t.Errorf("parseTakeCodeToMap[%s] newMapName[%s] mapName[%s]", v.data, newMapName, v.mapName)
			}
		}
	})
}

func Test_2(t *testing.T) {
	// 创建一个map[int]bool
	myMap := map[int]bool{
		1: true,
		2: true,
	}

	// 准备数据
	data := struct {
		MapName string
		MyMap   map[int]bool
	}{
		MapName: "myMap",
		MyMap:   myMap,
	}

	// 定义模板
	const tmpl = `var {{.MapName}} = {{printf "%#v" .MyMap}}`

	// 解析模板
	tmplParsed, err := template.New("example").Parse(tmpl)
	if err != nil {
		panic(err)
	}

	// 渲染模板到一个缓冲区
	var buf bytes.Buffer
	err = tmplParsed.Execute(&buf, data)
	if err != nil {
		panic(err)
	}

	// 输出渲染结果
	fmt.Println(buf.String())
}

package gomod

import (
	"bytes"
	"os"
	"strings"
)

var (
	spaceBytes = []byte(" ")
	lineBytes  = []byte("\n")
)

func GetGoModuleName(dir string) string {
	all, err := os.ReadFile(dir + "/go.mod")
	if err != nil {
		panic(err.Error())
	}

	start := bytes.Index(all, spaceBytes)
	if start == -1 {
		panic("not found space")
	}

	end := bytes.Index(all, lineBytes)
	if end == -1 {
		panic(`not found \n`)
	}

	return string(all[start+1 : end])
}

func GetGoModuleLastName(dir string) string {
	name := GetGoModuleName(dir)
	if pos := strings.LastIndex(name, "/"); pos != -1 {
		return name[pos+1:]
	}
	return name
}

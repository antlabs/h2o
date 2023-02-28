// Copyright 2021-2023 guonaihong. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package codemsg

import (
	"bytes"
	"fmt"
	"go/ast"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/antlabs/deepcopy"
)

func genCodeMsg(c *CodeMsg) {
	dir := ""

	args := c.Args

	types := strings.Split(c.TypeNames, ",")
	var tags []string

	g := Generator{
		trimPrefix:  c.Trimprefix,
		lineComment: c.Linecomment,
	}

	if len(args) == 1 && isDirectory(args[0]) {
		dir = args[0]
	} else {
		if len(tags) != 0 {
			log.Fatal("-tags option applies only to directories, not when files are specified")
		}
		dir = filepath.Dir(args[0])
	}

	g.parsePackage(args, tags)

	for _, typeName := range types {
		g.generateCodeMsg(c, typeName)
	}

	saveFile(c, dir, g.format(), "_codemsg.go", types[0])
	if c.Grpc {
		saveFile(c, dir, g.formatGrpc(), "_grpc_status.go", types[0])
	}
}

func saveFile(c *CodeMsg, dir string, src []byte, suffix string, types string) {

	outputName := c.Output
	if outputName == "" {
		baseName := fmt.Sprintf("%s%s", types, suffix)
		outputName = filepath.Join(dir, strings.ToLower(baseName))
	}

	err := ioutil.WriteFile(outputName, src, 0644)
	if err != nil {
		log.Fatalf("writing output: %s", err)
	}
}

// generate produces the String method for the named type.
func (g *Generator) generateCodeMsg(c *CodeMsg, typeName string) {

	values := make([]Value, 0, 100)

	for _, file := range g.pkg.files {

		// Set the state for this run of the walker.
		file.typeName = typeName
		file.values = nil
		if file.file != nil {
			ast.Inspect(file.file, file.genDecl)
			values = append(values, file.values...)
		}
	}

	if len(values) == 0 {
		log.Fatalf("no values defined for type %s", typeName)
	}

	tmpl := CodeMsgTmpl{AllVariable: values, TypeName: typeName}
	deepcopy.Copy(&tmpl, c).Do()
	tmpl.Args = strings.Join(os.Args[2:], " ")
	tmpl.PkgName = g.pkg.name

	//tmpl.Gen(os.Stdout)
	if err := tmpl.Gen(&g.buf); err != nil {
		io.Copy(os.Stdout, bytes.NewReader(g.buf.Bytes()))
	}

	if c.Grpc {

		grpcTmpl := GrpcCodeMsgTmpl{AllVariable: values, TypeName: typeName, StringMethod: c.StringMethod}
		deepcopy.Copy(&grpcTmpl, c).Do()
		grpcTmpl.Args = strings.Join(os.Args[2:], " ")
		grpcTmpl.PkgName = g.pkg.name

		if err := grpcTmpl.Gen(&g.grpcBuf); err != nil {
			io.Copy(os.Stdout, bytes.NewReader(g.grpcBuf.Bytes()))
		}
	}
}

package ast

import (
	"fmt"
	"go/ast"
	"go/doc"
	"go/parser"
	"go/token"
	"strings"

	"github.com/antlabs/h2o/internal/file"
)

type Ast struct {
	From string `clop:"short;long" usage:"which file to open from? If it is -, it means reading from stdin"`
	Dir  string `clop:"long" usage:"gen dir" default:""`
	Doc  bool   `clop:"short;long" usage:"Use the doc package to parse"`
}

func (a *Ast) doc() {

	fset := token.NewFileSet() // positions are relative to fset

	d, err := parser.ParseDir(fset, a.Dir, nil, parser.ParseComments)
	if err != nil {
		fmt.Println(err)
		return
	}

	for k, f := range d {
		fmt.Println("package", k)
		p := doc.New(f, a.Dir, 0)

		for _, t := range p.Types {
			fmt.Println("  type", t.Name)
			fmt.Println("    docs:", t.Doc)
			if strings.HasPrefix(t.Doc, "antlabs.valid") {

				for _, spec := range t.Decl.Specs {
					switch spec.(type) {
					case *ast.TypeSpec:
						typeSpec := spec.(*ast.TypeSpec)

						fmt.Printf("Struct: name=%s\n", typeSpec.Name.Name)

						switch typeSpec.Type.(type) {
						case *ast.StructType:
							structType := typeSpec.Type.(*ast.StructType)
							for _, field := range structType.Fields.List {
								i := field.Type.(*ast.Ident)
								fieldType := i.Name

								for _, name := range field.Names {
									fmt.Printf("\tField: name=%s type=%s, %s %s\n", name.Name, fieldType, field.Doc.Text(), field.Comment.Text())
								}

							}

						}
					}
				}
			}

		}
	}
}

func (a *Ast) SubMain() {
	if len(a.Dir) > 0 {
		a.doc()
		return
	}

	all, err := file.ReadFile(a.From)
	if err != nil {
		panic(err.Error())
	}

	fset := token.NewFileSet() // positions are relative to fset
	f, err := parser.ParseFile(fset, "", all, 0)
	if err != nil {
		panic(err)
	}

	// Print the AST.
	ast.Print(fset, f)
}

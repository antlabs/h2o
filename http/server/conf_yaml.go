package server

import (
	_ "embed"
	"io"
	"text/template"
)

//go:embed conf_yaml.tmpl
var httpConfigYamlTemplate string

type ConfigYAML struct {
	GoModLastName string
}

func (c *ConfigYAML) Gen(w io.Writer) {
	tpl := func() *template.Template {
		tmpl := httpConfigYamlTemplate
		return template.Must(template.New("h2o-http-server-config-yaml-tmpl").Parse(tmpl))
	}()

	tpl.Execute(w, c)
}

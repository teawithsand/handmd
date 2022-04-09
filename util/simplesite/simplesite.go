package simplesite

import (
	"html/template"
	"io"
)

const rawTemplate = `<!DOCTYPE html>
<html lang="{{ .Lang }}">
<head>
	<meta charset="utf-8">
	<title>{{ .Title }}</title>
	{{ range .Inputs }}
		{{ .Render}}
	{{ end }}
	{{range .Metas }}
		{{ .Render}}
	{{end}}
	{{range .Links}}
		{{ .Render }}
	{{end}}
	{{range .Scripts}}
		{{ .Render }}
	{{end}}
</head>
<body>
</body>
</html>`

func getCompiledTemplate() *template.Template {
	tpl, err := template.New("asdf").Parse(rawTemplate)
	if err != nil {
		panic(err)
	}

	return tpl
}

var compiledTemplate = getCompiledTemplate()

type RenderData struct {
	Head Head
}

func (srd *RenderData) Render(w io.Writer) (err error) {
	err = compiledTemplate.Execute(w, srd.Head)
	if err != nil {
		return
	}
	return
}

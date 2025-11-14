package handler

import (
	"html/template"
	"net/http"

	"github.com/f044fs3t5w3f/metrics/internal/repository"
)

var tmpl *template.Template = template.Must(template.New("template").Parse(`
<!DOCTYPE html>
<html>
  <head>
    <meta charset="utf-8">
    <title>Метрики</title>
  </head>
  <body>
	<h1>Метрики</h1>
	<li>
	{{range . }} <p>
		{{.MType}} {{.ID}}: {{ if (eq .MType "gauge") }} {{.Value}} {{else}} {{.Delta}} {{end}}
	</p> {{end}}
	</li>
  </body>
</html>
`))

func Index(storage repository.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		metrics, _ := storage.GetValuesList()
		tmpl.Execute(w, metrics)
	}
}

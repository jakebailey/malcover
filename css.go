package main

import (
	"html/template"
	"io"
	"log"
)

var cssTemplate = template.Must(template.New("css").Parse(`/* {{.Name}} */
#more{{.ID}} {
	background-image: url({{.URL}});
}
`))

func renderCSS(w io.Writer, name string, id int, url string) {
	err := cssTemplate.Execute(w, struct {
		ID   int
		Name string
		URL  string
	}{
		ID:   id,
		Name: name,
		URL:  url,
	})
	if err != nil {
		log.Println(err)
	}
}

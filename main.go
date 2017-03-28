package main

import (
	"io"
	"log"
	"net/http"
	"text/template"
	"time"

	"github.com/nstratos/go-myanimelist/mal"
	"github.com/pressly/chi"
	"github.com/pressly/chi/middleware"
	"github.com/tdewolff/minify"
	"github.com/tdewolff/minify/css"
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

var m = minify.New()

func init() {
	m.AddFunc("text/css", css.Minify)
}

func maybeMinify(w io.Writer, r *http.Request) io.Writer {
	if r.FormValue("minify") == "true" {
		return m.Writer("text/css", w)
	}
	return w
}

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.CloseNotify)
	r.Use(middleware.Heartbeat("/ping"))

	c := mal.NewClient(nil)

	r.Route("/:username", func(r chi.Router) {
		r.Use(middleware.ThrottleBacklog(5, 0, 5*time.Second))

		r.Get("/anime.css", func(w http.ResponseWriter, r *http.Request) {
			username := chi.URLParam(r, "username")

			list, _, err := c.Anime.List(username)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			if list.MyInfo.Name == "" {
				http.NotFound(w, r)
				return
			}

			w.Header().Set("Content-Type", "text/css; charset=utf-8")

			wr := maybeMinify(w, r)
			for _, v := range list.Anime {
				renderCSS(wr, v.SeriesTitle, v.SeriesAnimeDBID, v.SeriesImage)
			}
		})

		r.Get("/manga.css", func(w http.ResponseWriter, r *http.Request) {
			username := chi.URLParam(r, "username")

			list, _, err := c.Manga.List(username)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			if list.MyInfo.Name == "" {
				http.NotFound(w, r)
				return
			}

			w.Header().Set("Content-Type", "text/css; charset=utf-8")

			wr := maybeMinify(w, r)
			for _, v := range list.Manga {
				renderCSS(wr, v.SeriesTitle, v.SeriesMangaDBID, v.SeriesImage)
			}
		})
	})

	log.Fatal(http.ListenAndServe(":5000", r))
}

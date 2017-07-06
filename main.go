package main

import (
	"log"
	"net/http"
	"time"

	"github.com/alexflint/go-arg"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/nstratos/go-myanimelist/mal"
)

var args = struct {
	Port string
}{
	Port: "5000",
}

func main() {
	arg.MustParse(&args)

	log.SetFlags(log.LstdFlags | log.Lshortfile)

	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.CloseNotify)
	r.Use(middleware.Heartbeat("/ping"))

	c := mal.NewClient(nil)

	r.Route("/{username}", func(r chi.Router) {
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

	addr := ":" + args.Port
	log.Println("starting server at", addr)
	log.Fatal(http.ListenAndServe(addr, r))
}

package main

import (
	"log"
	"net/http"

	"github.com/tdewolff/minify"
	"github.com/tdewolff/minify/css"
)

func minifyMiddleware(next http.Handler) http.Handler {
	minifier := minify.New()
	minifier.AddFunc("text/css", css.Minify)

	fn := func(w http.ResponseWriter, r *http.Request) {
		if r.FormValue("minify") == "true" {
			mw := minifier.ResponseWriter(w, r)
			defer func() {
				if err := mw.Close(); err != nil {
					log.Println(err)
				}
			}()

			w = mw
		}

		next.ServeHTTP(w, r)
	}

	return http.HandlerFunc(fn)
}

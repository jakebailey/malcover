package main

import (
	"io"
	"net/http"

	"github.com/tdewolff/minify"
	"github.com/tdewolff/minify/css"
)

var minifier = minify.New()

func init() {
	minifier.AddFunc("text/css", css.Minify)
}

func maybeMinify(w io.Writer, r *http.Request) io.Writer {
	if r.FormValue("minify") == "true" {
		return minifier.Writer("text/css", w)
	}
	return w
}

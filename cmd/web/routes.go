package main

import (
	"net/http"

	"github.com/bmizerany/pat"
	"github.com/justinas/alice"
)

func (app *application) routes(cfg *Config) http.Handler {
	standardMiddleware := alice.New(app.recoverPanic, app.logRequest, secureHeaders)

	mux := pat.New()
	mux.Get("/", http.HandlerFunc(app.home))
	mux.Get("/snippets/create", http.HandlerFunc(app.createSnippetForm))
	mux.Post("/snippets/create", http.HandlerFunc(app.createSnippet))
	mux.Get("/snippets/:id", http.HandlerFunc(app.showSnippet))

	fileServer := http.FileServer(http.Dir(cfg.StaticDir))
	mux.Get("/static/", http.StripPrefix("/static", fileServer))

	return standardMiddleware.Then(mux)
}

package main

import "net/http"

func (app *application) routes(cfg *Config) *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/", app.home)
	mux.HandleFunc("/snippets", app.showSnippet)
	mux.HandleFunc("/snippets/create", app.createSnippet)

	fileServer := http.FileServer(http.Dir(cfg.StaticDir))
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	return mux
}

package main

import "net/http"

func (app *application) routes() http.Handler {
	// multiplexer and fileserver configuration
	mux := http.NewServeMux()
	fileServer := http.FileServer(http.Dir(app.staticDir))
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	// define routes
	mux.HandleFunc("/", app.home)
	mux.HandleFunc("/snippet/view", app.snippetView)
	mux.HandleFunc("/snippet/create", app.snippetCreate)
	return app.logRequest(secureHeaders(mux))
}

package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/justinas/alice"
)

func (app *application) routes() http.Handler {
	// multiplexer and fileserver configuration
	router := httprouter.New()
	router.NotFound = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		app.notFound(w)
	})

	fileServer := http.FileServer(http.Dir(app.staticDir))
	router.Handler(http.MethodGet, "/static/*filepath", http.StripPrefix("/static", fileServer))

	dynamic := alice.New(app.sessionManager.LoadAndSave)

	// define routes
	router.Handler(http.MethodGet, "/", dynamic.ThenFunc(app.home))
	router.Handler(http.MethodGet, "/snippet/view/:id", dynamic.ThenFunc(app.snippetView))
	router.Handler(http.MethodGet, "/snippet/create", dynamic.ThenFunc(app.snippetCreate))
	router.Handler(http.MethodPost, "/snippet/create", dynamic.ThenFunc(app.snippetCreatePost))

	// define routes for user authentication
	router.Handler(http.MethodGet, "/user/signup", dynamic.ThenFunc(app.userSignup))
	router.Handler(http.MethodGet, "/user/login", dynamic.ThenFunc(app.userLogin))
	router.Handler(http.MethodPost, "/user/login", dynamic.ThenFunc(app.userLoginPost))
	router.Handler(http.MethodPost, "/user/logout", dynamic.ThenFunc(app.userLogoutPost))
	router.Handler(http.MethodPost, "/user/signup", dynamic.ThenFunc(app.userSignupPost))

	// middleware chaining
	mdlw := alice.New(app.recoverPanic, app.logRequest, secureHeaders)

	return mdlw.Then(router)
}

package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/justinas/alice"
)

func (app *application) routes() http.Handler {
	router := httprouter.New()

	// use our custom notFound() helper to handle 404 in all cases
	router.NotFound = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		app.notFound(w)
	})

	fileServer := http.FileServer(http.Dir("./ui/static/"))
	router.Handler(http.MethodGet, "/static/*filepath", http.StripPrefix("/static", fileServer))

	// create a new middleware chain containing middleware
	// specific to dynamic routes
	dynamicRoutesMiddleware := alice.New(app.sessionManager.LoadAndSave)

	router.Handler(http.MethodGet, "/", dynamicRoutesMiddleware.ThenFunc(app.home))
	router.Handler(http.MethodGet, "/snippet/view/:id", dynamicRoutesMiddleware.ThenFunc(app.snippetView))
	router.Handler(http.MethodGet, "/snippet/create", dynamicRoutesMiddleware.ThenFunc(app.snippetCreate))
	router.Handler(http.MethodPost, "/snippet/create", dynamicRoutesMiddleware.ThenFunc(app.snippetCreatePost))

	// create a middleware chain using 3rd party package
	standardMiddleware := alice.New(app.recoverPanic, app.logRequest, secureHeaders)
	return standardMiddleware.Then(router)
}

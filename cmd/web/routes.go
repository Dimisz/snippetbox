package main

import (
	"net/http"
	"snippetbox/ui"

	"github.com/justinas/alice"
)

func (app *application) routes() http.Handler {
	mux := http.NewServeMux()

	// serving static assets from disk
	// fileServer := http.FileServer(http.Dir("./ui/static"))
	// mux.Handle("GET /static/", http.StripPrefix("/static", fileServer))

	// embedding static files hence compiled to a binary
	mux.Handle("GET /static/", http.FileServerFS(ui.Files))

	// middleware for session management
	dynamic := alice.New(app.sessionManager.LoadAndSave, noSurf, app.authenticate)

	// unprotected routes, i.e. do not require auth
	// using only 'dynamic' middleware chain
	mux.Handle("GET /{$}", dynamic.ThenFunc(app.home))
	mux.Handle("GET /snippet/view/{id}", dynamic.ThenFunc(app.snippetView))
	mux.Handle("GET /user/signup", dynamic.ThenFunc(app.userSignup))
	mux.Handle("POST /user/signup", dynamic.ThenFunc(app.userSignupPost))
	mux.Handle("GET /user/login", dynamic.ThenFunc(app.userLogin))
	mux.Handle("POST /user/login", dynamic.ThenFunc(app.userLoginPost))

	// protected routes, i.e. require auth
	// use 'protected' middleware chain, i.e. dynamic + requireAuth
	protected := dynamic.Append(app.requireAuthentication)

	mux.Handle("GET /snippet/create", protected.ThenFunc(app.snippetCreate))
	mux.Handle("POST /snippet/create", protected.ThenFunc(app.snippetCreatePost))
	mux.Handle("POST /user/logout", protected.ThenFunc(app.userLogoutPost))

	standard := alice.New(app.recoverPanic, app.logRequest, commonHeaders)

	return standard.Then(mux)
}

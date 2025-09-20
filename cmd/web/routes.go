package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/justinas/alice"
)

func (app *application) routes() http.Handler {
	router := httprouter.New()
	// override default NotFound method to return our custom app.notFound()
	router.NotFound = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		app.notFound(w)
	})

	fileServer := http.FileServer(http.Dir("./ui/static/"))
	router.Handler(http.MethodGet, "/static/*filepath", http.StripPrefix("/static", fileServer))

	dynamicMiddlewareChain := alice.New(app.sessionManager.LoadAndSave, noSurf)

	// auth not required
	router.Handler(http.MethodGet, "/", dynamicMiddlewareChain.ThenFunc(app.home))
	router.Handler(http.MethodGet, "/snippet/view/:id", dynamicMiddlewareChain.ThenFunc(app.snippetView))
	router.Handler(http.MethodGet, "/user/signup", dynamicMiddlewareChain.ThenFunc(app.userSignup))
	router.Handler(http.MethodPost, "/user/signup", dynamicMiddlewareChain.ThenFunc(app.userSignupPost))
	router.Handler(http.MethodGet, "/user/login", dynamicMiddlewareChain.ThenFunc(app.userLogin))
	router.Handler(http.MethodPost, "/user/login", dynamicMiddlewareChain.ThenFunc(app.userLoginPost))

	// auth required
	protected := dynamicMiddlewareChain.Append(app.requireAuthentication)

	router.Handler(http.MethodGet, "/snippet/create", protected.ThenFunc(app.snippetCreate))
	router.Handler(http.MethodPost, "/snippet/create", protected.ThenFunc(app.snippetCreatePost))
	router.Handler(http.MethodPost, "/user/logout", protected.ThenFunc(app.userLogoutPost))

	standardMiddlewareChain := alice.New(app.recoverPanic, app.logRequest, secureHeaders)
	return standardMiddlewareChain.Then(router)
}

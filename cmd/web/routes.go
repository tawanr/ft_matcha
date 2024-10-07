package main

import (
	"net/http"

	"github.com/justinas/alice"
	"github.org/tawanr/ft_matcha/ui"
)

func (app *application) routes() http.Handler {
	mux := http.NewServeMux()

	mux.Handle("GET /static/", http.FileServerFS(ui.Files))
	mux.HandleFunc("GET /ping", ping)

	fileServer := http.FileServer(http.Dir("./ui/public/"))
	mux.Handle("GET /public/", http.StripPrefix("/public/", fileServer))

	dynamic := alice.New(app.sessionManager.LoadAndSave, noSurf, app.authenticate)

	mux.Handle("GET /{$}", dynamic.ThenFunc(app.home))
	mux.Handle("GET /user", dynamic.ThenFunc(app.userList))
	mux.Handle("GET /user/signup", dynamic.ThenFunc(app.userSignup))
	mux.Handle("POST /user/signup", dynamic.ThenFunc(app.userSignupPost))
	mux.Handle("GET /user/login", dynamic.ThenFunc(app.userLogin))
	mux.Handle("POST /user/login", dynamic.ThenFunc(app.userLoginPost))

	protected := dynamic.Append(app.requireAuthentication)
	mux.Handle("POST /user/logout", protected.ThenFunc(app.userLogoutPost))
	mux.Handle("GET /user/profile", protected.ThenFunc(app.userProfile))
	mux.Handle("POST /user/profile", protected.ThenFunc(app.userProfilePost))

	standard := alice.New(app.recoverPanic, app.logRequest, commonHeaders)
	return standard.Then(mux)
}

package main

import "net/http"

func (app *application) routes() http.Handler {
	mux := http.NewServeMux()

	fileServer := http.FileServer(http.Dir("./ui/static/"))
	mux.Handle("GET /static/", http.StripPrefix("/static", fileServer))

	mux.HandleFunc("GET /{$}", app.home)
	mux.HandleFunc("GET /user", app.userList)
	mux.HandleFunc("GET /user/signup", app.userSignup)
	mux.HandleFunc("POST /user/signup", app.userSignupPost)
	return app.recoverPanic(app.logRequest(commonHeaders(mux)))
}

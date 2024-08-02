package main

import (
	"fmt"
	"html/template"
	"net/http"
)

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	ts, err := template.ParseFiles("./ui/html/pages/home1.tmpl")
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	err = ts.Execute(w, nil)
	if err != nil {
		app.serverError(w, r, err)
	}
}

func (app *application) userList(w http.ResponseWriter, r *http.Request) {

}

func (app *application) userSignup(w http.ResponseWriter, r *http.Request) {
	ts, err := template.ParseFiles("./ui/html/pages/user/signup.tmpl")
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	err = ts.Execute(w, nil)
	if err != nil {
		app.serverError(w, r, err)
	}
}

func (app *application) userSignupPost(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "signup post")
}

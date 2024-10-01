package main

import (
	"fmt"
	"net/http"

	"github.org/tawanr/ft_matcha/internal/validator"
)

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)
	app.render(w, r, http.StatusOK, "home.go.tmpl", data)
}

func (app *application) userList(w http.ResponseWriter, r *http.Request) {

}

type userSignupForm struct {
	Name                string `form:"name"`
	Email               string `form:"email"`
	Password            string `form:"password"`
	validator.Validator `form:"-"`
}

func (app *application) userSignup(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)
	data.Form = userSignupForm{}
	app.render(w, r, http.StatusOK, "signup.go.tmpl", data)
}

func (app *application) userSignupPost(w http.ResponseWriter, r *http.Request) {
	var form userSignupForm

	err := app.decodePostForm(r, &form)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	form.CheckField(validator.NotBlank("name"), "name", "Name is required")
	form.CheckField(validator.NotBlank("email"), "email", "Email is required")
	form.CheckField(validator.NotBlank("password"), "password", "Password is required")
	form.CheckField(validator.MinChars("password", 8), "password", "Password must be at least 8 characters")

	if !form.Valid() {
		data := app.newTemplateData(r)
		data.Form = form
		app.render(w, r, http.StatusOK, "signup.go.tmpl", data)
		return
	}

	fmt.Printf("signup post %+v\n", form)
	app.sessionManager.Put(r.Context(), "registered", form.Email)

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

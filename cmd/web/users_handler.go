package main

import (
	"errors"
	"fmt"
	"net/http"

	"github.org/tawanr/ft_matcha/internal/models"
	"github.org/tawanr/ft_matcha/internal/validator"
)

func (app *application) userList(w http.ResponseWriter, r *http.Request) {

}

type userSignupForm struct {
	Name                string `form:"name"`
	Email               string `form:"email"`
	Password            string `form:"password"`
	validator.Validator `form:"-"`
}

type userLoginForm struct {
	Email               string `form:"email"`
	Password            string `form:"password"`
	validator.Validator `form:"-"`
}

type profileForm struct {
	Bio                 string `form:"bio"`
	Gender              int    `form:"gender"`
	Preferences         []int  `form:"preferences"`
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

	fmt.Println("name", validator.NotBlank(form.Name))
	form.CheckField(validator.NotBlank(form.Name), "name", "Name is required")
	form.CheckField(validator.NotBlank(form.Email), "email", "Email is required")
	form.CheckField(validator.Matches(form.Email, validator.EmailRX), "email", "Email is invalid")
	form.CheckField(validator.NotBlank(form.Password), "password", "Password is required")
	form.CheckField(validator.MinChars(form.Password, 8), "password", "Password must be at least 8 characters")

	if !form.Valid() {
		data := app.newTemplateData(r)
		data.Form = form
		app.render(w, r, http.StatusOK, "signup.go.tmpl", data)
		return
	}

	err = app.users.Insert(form.Name, form.Email, form.Password)
	if err != nil {
		if errors.Is(err, models.ErrDuplicateEmail) {
			form.AddFieldError("email", "Email already exists")
			data := app.newTemplateData(r)
			data.Form = form
			app.render(w, r, http.StatusOK, "signup.go.tmpl", data)
		} else {
			app.serverError(w, r, err)
		}
		return
	}
	app.sessionManager.Put(r.Context(), "registered", form.Email)

	http.Redirect(w, r, "/user/login", http.StatusSeeOther)
}

func (app *application) userLogin(w http.ResponseWriter, r *http.Request) {
	if app.sessionManager.GetInt(r.Context(), "authenticatedUserID") != 0 {
		http.Redirect(w, r, "/user/profile", http.StatusSeeOther)
		return
	}
	data := app.newTemplateData(r)
	data.Form = userLoginForm{}
	app.render(w, r, http.StatusOK, "login.go.tmpl", data)
}

func (app *application) userLoginPost(w http.ResponseWriter, r *http.Request) {
	var form userLoginForm

	err := app.decodePostForm(r, &form)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	form.CheckField(validator.NotBlank(form.Email), "email", "Email is required")
	form.CheckField(validator.Matches(form.Email, validator.EmailRX), "email", "Email is invalid")
	form.CheckField(validator.NotBlank(form.Password), "password", "Password is required")

	if !form.Valid() {
		data := app.newTemplateData(r)
		data.Form = form
		app.render(w, r, http.StatusOK, "login.go.tmpl", data)
		return
	}

	id, err := app.users.Authenticate(form.Email, form.Password)
	if err != nil {
		if errors.Is(err, models.ErrInvalidCredentials) {
			form.AddNonFieldError("Invalid email or password")
			data := app.newTemplateData(r)
			data.Form = form
			app.render(w, r, http.StatusOK, "login.go.tmpl", data)
		} else {
			app.serverError(w, r, err)
		}
		return
	}

	err = app.sessionManager.RenewToken(r.Context())
	if err != nil {
		app.serverError(w, r, err)
		return
	}
	app.sessionManager.Put(r.Context(), "authenticatedUserID", id)
	http.Redirect(w, r, "/user/profile", http.StatusSeeOther)
}

func (app *application) userLogoutPost(w http.ResponseWriter, r *http.Request) {
	err := app.sessionManager.RenewToken(r.Context())
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	app.sessionManager.Remove(r.Context(), "authenticatedUserID")
	app.sessionManager.Put(r.Context(), "flash", "You've been logged out")

	http.Redirect(w, r, "/user/login", http.StatusSeeOther)
}

func (app *application) userProfile(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)
	data.Form = profileForm{}
	app.render(w, r, http.StatusOK, "profile.go.tmpl", data)
}

func (app *application) userProfilePost(w http.ResponseWriter, r *http.Request) {
	var form profileForm

	err := app.decodePostForm(r, &form)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	if !form.Valid() {
		data := app.newTemplateData(r)
		data.Form = form
		app.render(w, r, http.StatusOK, "profile.go.tmpl", data)
		return
	}

	fmt.Printf("%#v\n", form)

	// err = app.users.UpdateProfile(r.Context(), form.Gender, form.Preferences)
	// if err != nil {
	// 	if errors.Is(err, models.ErrNoRecord) {
	// 		app.clientError(w, http.StatusBadRequest)
	// 		return
	// 	}
	// 	app.serverError(w, r, err)
	// 	return
	// }

	http.Redirect(w, r, "/user/profile", http.StatusSeeOther)
}

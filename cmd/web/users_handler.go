package main

import (
	"cmp"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/tawanr/ft_matcha/internal/models"
	"github.com/tawanr/ft_matcha/internal/validator"
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
	Bio                 string            `form:"bio"`
	Age                 int               `form:"age"`
	Gender              models.GenderType `form:"gender"`
	Preferences         []int             `form:"preferences"`
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

	err = app.models.Users.Insert(form.Name, form.Email, form.Password)
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

	id, err := app.models.Users.Authenticate(form.Email, form.Password)
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
	userID := app.sessionManager.GetInt(r.Context(), "authenticatedUserID")
	if userID == 0 {
		http.Redirect(w, r, "/user/login", http.StatusSeeOther)
		return
	}
	profile, err := app.models.Profiles.Get(userID)
	if err != nil {
		if !errors.Is(err, models.ErrNoRecord) {
			app.serverError(w, r, err)
			return
		}
		profile, err = models.CreateDefaultProfile(userID, app.models.Profiles)
		if err != nil {
			app.serverError(w, r, err)
			return
		}
	}
	data := app.newTemplateData(r)
	preferences := []int{}
	if profile.PreferMale {
		preferences = append(preferences, 0)
	}
	if profile.PreferFemale {
		preferences = append(preferences, 1)
	}
	form := profileForm{}
	form.Bio = profile.Bio
	form.Gender = profile.Gender
	form.Preferences = preferences
	form.Age = profile.Age
	data.Form = form
	data.Profile = *profile
	fmt.Println(data.Profile)
	app.render(w, r, http.StatusOK, "profile.go.tmpl", data)
}

func (app *application) userProfilePost(w http.ResponseWriter, r *http.Request) {
	userID := app.sessionManager.GetInt(r.Context(), "authenticatedUserID")
	if userID == 0 {
		http.Redirect(w, r, "/user/login", http.StatusSeeOther)
		return
	}
	var form profileForm

	err := app.decodePostForm(r, &form)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	form.CheckField(cmp.Compare(form.Age, 18) >= 0, "age", "Age must be at least 18")

	if !form.Valid() {
		data := app.newTemplateData(r)
		data.Form = form
		app.render(w, r, http.StatusOK, "profile.go.tmpl", data)
		return
	}

	_, err = app.models.Profiles.Update(userID, form.Gender, form.Preferences, form.Bio, form.Age)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			app.clientError(w, http.StatusBadRequest)
			return
		}
		app.serverError(w, r, err)
		return
	}

	http.Redirect(w, r, "/user/profile", http.StatusSeeOther)
}

func (app *application) userProfileCard(w http.ResponseWriter, r *http.Request) {
	userID, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		app.clientError(w, http.StatusNotFound)
		return
	}
	if userID == 0 {
		userID = app.sessionManager.GetInt(r.Context(), "authenticatedUserID")
		if userID == 0 {
			app.clientError(w, http.StatusBadRequest)
			return
		}
	}

	profile, err := app.models.Profiles.Get(userID)
	if err != nil {
		if !errors.Is(err, models.ErrNoRecord) {
			app.serverError(w, r, err)
			return
		}
		profile, err = models.CreateDefaultProfile(userID, app.models.Profiles)
		if err != nil {
			app.serverError(w, r, err)
			return
		}
	}
	data := app.newTemplateData(r)
	data.Profile = *profile
	app.renderPartial(w, r, http.StatusOK, "user/profile-card.go.tmpl", data)
}

func (app *application) userProfileImage(w http.ResponseWriter, r *http.Request) {
	userID := app.sessionManager.GetInt(r.Context(), "authenticatedUserID")
	if userID == 0 {
		app.clientError(w, http.StatusBadRequest)
		return
	}
	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	file, handler, err := r.FormFile("user-image")
	if err != nil {
		app.serverError(w, r, err)
		return
	}
	defer file.Close()

	contentType := handler.Header.Get("Content-Type")
	if contentType != "image/png" && contentType != "image/jpeg" {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	timestamp := strconv.FormatInt(time.Now().Unix(), 10)
	newFilename := fmt.Sprintf("profile-%s-%s", timestamp, handler.Filename)

	uploadDir := "ui/public/images"
	publicDir := "/public/images"
	filePath := filepath.Join(uploadDir, newFilename)
	publicPath := filepath.Join(publicDir, newFilename)
	if _, err := os.Stat(uploadDir); os.IsNotExist(err) {
		err = os.MkdirAll(uploadDir, os.ModePerm)
		if err != nil {
			app.serverError(w, r, err)
			return
		}
	}

	dst, err := os.Create(filePath)
	if err != nil {
		app.serverError(w, r, err)
		return
	}
	defer dst.Close()

	_, err = io.Copy(dst, file)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	fmt.Printf("File Saved: %s\n", filePath)

	err = app.models.Profiles.AddImage(userID, models.Image(publicPath))
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	app.userProfileCard(w, r)
}

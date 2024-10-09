package main

import (
	"encoding/json"
	"html/template"
	"io/fs"
	"path/filepath"
	"time"

	"github.com/tawanr/ft_matcha/internal/models"
	"github.com/tawanr/ft_matcha/ui"
)

type templateData struct {
	CurrentYear     int
	User            models.User
	Form            any
	Flash           string
	isAuthenticated bool
	CSRFToken       string
	Profile         models.Profile
}

func humanDate(t time.Time) string {
	if t.IsZero() {
		return ""
	}
	return t.Format("2 Jan 2006")
}

func hasValue[T comparable](values []T, value T) bool {
	for _, v := range values {
		if v == value {
			return true
		}
	}
	return false
}

func genderString(g models.GenderType) string {
	switch g {
	case 0:
		return "Male"
	case 1:
		return "Female"
	case 2:
		return "Other"
	default:
		return ""
	}
}

func jsonMarshal(v any) string {
	b, err := json.Marshal(v)
	if err != nil {
		return ""
	}
	return string(b)
}

var functions = template.FuncMap{
	"humanDate":    humanDate,
	"hasValueInt":  hasValue[int],
	"genderString": genderString,
	"jsonMarshal":  jsonMarshal,
}

func newTemplateCache() (map[string]*template.Template, error) {
	cache := map[string]*template.Template{}

	pages, err := fs.Glob(ui.Files, "html/pages/*.go.tmpl")
	if err != nil {
		return nil, err
	}
	nested, err := fs.Glob(ui.Files, "html/pages/**/*.go.tmpl")
	pages = append(pages, nested...)

	if err != nil {
		return nil, err
	}

	for _, page := range pages {
		name := filepath.Base(page)
		patterns := []string{
			"html/base.go.tmpl",
			"html/partials/*.go.tmpl",
			page,
		}
		ts, err := template.New(name).Funcs(functions).ParseFS(ui.Files, patterns...)
		if err != nil {
			return nil, err
		}

		cache[name] = ts
	}
	return cache, nil
}

package main

import (
	"fmt"
	"html/template"
	"net/http"
)

func (app *application) Home(w http.ResponseWriter, r *http.Request) {
	partials := []string{
		"cmd/web/templates/base.layout.gohtml",
		"cmd/web/templates/header.partial.gohtml",
		"cmd/web/templates/footer.partial.gohtml",
		"cmd/web/templates/navbar.partial.gohtml",
		"cmd/web/templates/service-info.partial.gohtml",
	}

	var templateSlice []string
	templateSlice = append(templateSlice, fmt.Sprintf("cmd/web/templates/%s", "home.page.gohtml"))
	templateSlice = append(templateSlice, partials...)

	tmpl, err := template.ParseFiles(templateSlice...)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	td := &templateData{API: app.config.api}

	if err := tmpl.Execute(w, td); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (app *application) LoginPage(w http.ResponseWriter, r *http.Request) {
	partials := []string{
		"cmd/web/templates/base.layout.gohtml",
		"cmd/web/templates/header.partial.gohtml",
		"cmd/web/templates/footer.partial.gohtml",
		"cmd/web/templates/navbar.partial.gohtml",
		"cmd/web/templates/service-info.partial.gohtml",
	}

	var templateSlice []string
	templateSlice = append(templateSlice, fmt.Sprintf("cmd/web/templates/%s", "login.page.gohtml"))
	templateSlice = append(templateSlice, partials...)

	tmpl, err := template.ParseFiles(templateSlice...)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	td := &templateData{API: app.config.api}

	if err := tmpl.Execute(w, td); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (app *application) AllUsers(w http.ResponseWriter, r *http.Request) {
	partials := []string{
		"cmd/web/templates/base.layout.gohtml",
		"cmd/web/templates/header.partial.gohtml",
		"cmd/web/templates/footer.partial.gohtml",
		"cmd/web/templates/navbar.partial.gohtml",
		"cmd/web/templates/service-info.partial.gohtml",
	}

	var templateSlice []string
	templateSlice = append(templateSlice, fmt.Sprintf("cmd/web/templates/%s", "all-users.page.gohtml"))
	templateSlice = append(templateSlice, partials...)

	tmpl, err := template.ParseFiles(templateSlice...)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	td := &templateData{API: app.config.api}

	if err := tmpl.Execute(w, td); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

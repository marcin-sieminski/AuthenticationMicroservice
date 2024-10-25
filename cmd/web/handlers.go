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
	}

	var templateSlice []string
	templateSlice = append(templateSlice, fmt.Sprintf("cmd/web/templates/%s", "home.page.gohtml"))
	templateSlice = append(templateSlice, partials...)

	tmpl, err := template.ParseFiles(templateSlice...)
	if err != nil {
		app.errorLog.Println(err)
		return
	}

	td := app.addDefaultData(r)

	if err := tmpl.Execute(w, td); err != nil {
		app.errorLog.Println(err)
	}
}

func (app *application) LoginPage(w http.ResponseWriter, r *http.Request) {
	partials := []string{
		"cmd/web/templates/base.layout.gohtml",
		"cmd/web/templates/header.partial.gohtml",
		"cmd/web/templates/footer.partial.gohtml",
		"cmd/web/templates/navbar.partial.gohtml",
	}

	var templateSlice []string
	templateSlice = append(templateSlice, fmt.Sprintf("cmd/web/templates/%s", "login.page.gohtml"))
	templateSlice = append(templateSlice, partials...)

	tmpl, err := template.ParseFiles(templateSlice...)
	if err != nil {
		app.errorLog.Println(err)
		return
	}

	td := app.addDefaultData(r)

	if err := tmpl.Execute(w, td); err != nil {
		app.errorLog.Println(err)
	}
}

func (app *application) PostLoginPage(w http.ResponseWriter, r *http.Request) {
	app.Session.RenewToken(r.Context())

	err := r.ParseForm()
	if err != nil {
		return
	}

	email := r.Form.Get("email")
	password := r.Form.Get("password")

	id, err := app.DB.Authenticate(email, password)
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	user, err := app.DB.GetOneUser(id)
	if err != nil {
		return
	}
	app.Session.Put(r.Context(), "userName", fmt.Sprintf("%s %s", user.FirstName, user.LastName))
	app.Session.Put(r.Context(), "userID", id)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (app *application) Logout(w http.ResponseWriter, r *http.Request) {
	app.Session.Destroy(r.Context())
	app.Session.RenewToken(r.Context())

	http.Redirect(w, r, "/login", http.StatusSeeOther)
}

func (app *application) AllUsers(w http.ResponseWriter, r *http.Request) {
	partials := []string{
		"cmd/web/templates/base.layout.gohtml",
		"cmd/web/templates/header.partial.gohtml",
		"cmd/web/templates/footer.partial.gohtml",
		"cmd/web/templates/navbar.partial.gohtml",
	}

	var templateSlice []string
	templateSlice = append(templateSlice, fmt.Sprintf("cmd/web/templates/%s", "all-users.page.gohtml"))
	templateSlice = append(templateSlice, partials...)

	tmpl, err := template.ParseFiles(templateSlice...)
	if err != nil {
		app.errorLog.Println(err)
		return
	}

	td := app.addDefaultData(r)

	if err := tmpl.Execute(w, td); err != nil {
		app.errorLog.Println(err)
	}
}

func (app *application) OneUser(w http.ResponseWriter, r *http.Request) {
	partials := []string{
		"cmd/web/templates/base.layout.gohtml",
		"cmd/web/templates/header.partial.gohtml",
		"cmd/web/templates/footer.partial.gohtml",
		"cmd/web/templates/navbar.partial.gohtml",
	}

	var templateSlice []string
	templateSlice = append(templateSlice, fmt.Sprintf("cmd/web/templates/%s", "one-user.page.gohtml"))
	templateSlice = append(templateSlice, partials...)

	tmpl, err := template.ParseFiles(templateSlice...)
	if err != nil {
		app.errorLog.Println(err)
		return
	}

	td := app.addDefaultData(r)

	if err := tmpl.Execute(w, td); err != nil {
		app.errorLog.Println(err)
	}
}

func (app *application) addDefaultData(r *http.Request) *templateData {
	td := &templateData{}
	td.API = app.config.api

	if app.Session.Exists(r.Context(), "userID") {
		td.IsAuthenticated = true
		td.UserID = app.Session.GetInt(r.Context(), "userID")
		td.UserName = app.Session.GetString(r.Context(), "userName")
	} else {
		td.IsAuthenticated = false
		td.UserID = 0
	}

	return td
}

package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func (app *application) routes() http.Handler {
	mux := chi.NewRouter()
	mux.Use(SessionLoad)

	fileServer := http.FileServer(http.Dir("./static"))
	mux.Handle("/static/*", http.StripPrefix("/static", fileServer))

	mux.Get("/", app.Home)
	mux.Get("/ws", app.WsEndPoint)

	mux.Get("/login", app.LoginPage)
	mux.Post("/login", app.PostLoginPage)
	mux.Get("/logout", app.Logout)

	mux.Route("/admin", func(mux chi.Router) {
		mux.Use(app.Auth)
		mux.Get("/all-users", app.AllUsers)
		mux.Get("/all-users/{id}", app.OneUser)
	})

	return mux
}

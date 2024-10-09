package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func (app *application) routes() http.Handler {
	mux := chi.NewRouter()

	mux.Get("/login", app.LoginPage)

	mux.Route("/admin", func(mux chi.Router) {
		mux.Use(app.Auth)
	})

	return mux
}

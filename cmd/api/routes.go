package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
)

func (app *application) routes() http.Handler {
	mux := chi.NewRouter()

	mux.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	mux.Post("/authenticate", app.Authenticate)
	mux.Post("/is-authenticated", app.CheckAuthentication)

	mux.Route("/auth", func(mux chi.Router) {
		mux.Use(app.Auth)

		mux.Post("/all-users", app.AllUsers)
		mux.Post("/all-users/{id}", app.OneUser)
		mux.Post("/all-users/edit/{id}", app.EditUser)
		mux.Post("/all-users/delete/{id}", app.DeleteUser)
		mux.Post("/ask", app.Ask)
	})

	return mux
}

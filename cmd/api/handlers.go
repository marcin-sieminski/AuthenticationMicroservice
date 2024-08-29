package main

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/marcin-sieminski/AuthenticationService/models"
)

func (app *application) Authenticate(w http.ResponseWriter, r *http.Request) {
	var requestPayload struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	err := app.readJSON(w, r, &requestPayload)
	if err != nil {
		app.errorJSON(w, err, http.StatusBadRequest)
		return
	}

	user, err := app.DB.GetUserByEmail(requestPayload.Email)
	if err != nil {
		app.errorJSON(w, errors.New("błędne dane logowania"), http.StatusBadRequest)
		return
	}

	valid, err := app.DB.PasswordMatches(user.Password, requestPayload.Password)
	if err != nil || !valid {
		app.invalidCredentials(w)
		return
	}

	token, err := models.GenerateToken(user.ID, 24*time.Hour, models.ScopeAuthentication)
	if err != nil {
		app.invalidCredentials(w)
		return
	}

	err = app.DB.InsertToken(token, user)
	if err != nil {
		app.badRequest(w, r, err)
		return
	}

	payload := jsonResponse{
		Error:   false,
		Message: fmt.Sprintf("Zalogowano użytkownika %s", user.Email),
		User:    user,
		Token:   token,
	}

	app.writeJSON(w, http.StatusAccepted, payload)
}

func (app *application) CheckAuthentication(w http.ResponseWriter, r *http.Request) {
	user, err := app.authenticateToken(r)
	if err != nil {
		app.invalidCredentials(w)
		return
	}

	var payload struct {
		Error   bool   `json:"error"`
		Message string `json:"message"`
	}
	payload.Error = false
	payload.Message = fmt.Sprintf("Zalogowano użytkownika %s", user.Email)
	app.writeJSON(w, http.StatusOK, payload)
}

func (app *application) authenticateToken(r *http.Request) (*models.User, error) {
	authorizationHeader := r.Header.Get("Authorization")
	if authorizationHeader == "" {
		return nil, errors.New("niepoprawne dane autoryzacyjne")
	}

	headerParts := strings.Split(authorizationHeader, " ")
	if len(headerParts) != 2 || headerParts[0] != "Bearer" {
		return nil, errors.New("niepoprawne dane autoryzacyjne")
	}

	token := headerParts[1]
	if len(token) != 26 {
		return nil, errors.New("niepoprawny token autoryzacyjny")
	}

	user, err := app.DB.GetUserForToken(token)
	if err != nil {
		return nil, errors.New("nie znaleziono użytkownika")
	}

	return user, nil
}

package main

import (
	"errors"
	"fmt"
	"net/http"
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

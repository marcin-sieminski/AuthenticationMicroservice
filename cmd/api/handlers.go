package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/marcin-sieminski/AuthenticationService/models"
	"github.com/tmc/langchaingo/llms/openai"
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
	payload.Message = fmt.Sprintf("Użytkownik %s zalogowany", user.Email)
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

func (app *application) AllUsers(w http.ResponseWriter, r *http.Request) {
	allUsers, err := app.DB.GetAllUsers()
	if err != nil {
		app.badRequest(w, r, err)
		return
	}

	app.writeJSON(w, http.StatusOK, allUsers)
}

func (app *application) OneUser(w http.ResponseWriter, r *http.Request) {
	id := strings.Split(r.URL.Path, "/")[3]
	userID, _ := strconv.Atoi(id)

	user, err := app.DB.GetOneUser(userID)
	if err != nil {
		app.badRequest(w, r, err)
		return
	}

	app.writeJSON(w, http.StatusOK, user)
}

func (app *application) EditUser(w http.ResponseWriter, r *http.Request) {
	id := strings.Split(r.URL.Path, "/")[4]
	userID, _ := strconv.Atoi(id)

	var user models.User

	err := app.readJSON(w, r, &user)
	if err != nil {
		app.badRequest(w, r, err)
		return
	}

	if userID > 0 {
		err = app.DB.Update(user)
		if err != nil {
			app.badRequest(w, r, err)
			return
		}

		if user.Password != "" {
			err = app.DB.ResetPassword(user.Password, userID)
			if err != nil {
				app.badRequest(w, r, err)
				return
			}
		}
	} else {
		_, err = app.DB.Insert(user)
		if err != nil {
			app.badRequest(w, r, err)
			return
		}
	}

	var resp struct {
		Error   bool   `json:"error"`
		Message string `json:"message"`
	}

	resp.Error = false
	app.writeJSON(w, http.StatusOK, resp)
}

func (app *application) DeleteUser(w http.ResponseWriter, r *http.Request) {
	id := strings.Split(r.URL.Path, "/")[4]
	userID, _ := strconv.Atoi(id)

	err := app.DB.Delete(userID)
	if err != nil {
		app.badRequest(w, r, err)
		return
	}

	var resp struct {
		Error   bool   `json:"error"`
		Message string `json:"message"`
	}

	resp.Error = false
	app.writeJSON(w, http.StatusOK, resp)
}

func (app *application) Ask(w http.ResponseWriter, r *http.Request) {
	var requestData struct {
		Prompt string `json:"prompt"`
	}
	if err := app.readJSON(w, r, &requestData); err != nil {
		app.badRequest(w, r, err)
		return
	}

	llm, err := openai.New()
	if err != nil {
		app.badRequest(w, r, err)
		return
	}

	completion, err := llm.Call(context.Background(), requestData.Prompt)
	if err != nil {
		app.badRequest(w, r, err)
		return
	}

	app.writeJSON(w, http.StatusOK, completion)
}

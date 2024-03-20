package main

import (
	"errors"
	"fmt"
	"net/http"
)

func (app *Config) Authenticate(w http.ResponseWriter, r *http.Request) {
	var request struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	err := app.ReadJSON(w, r, &request)
	if err != nil {
		app.ErrorJSON(w, err, http.StatusBadRequest)
		return
	}

	user, err := app.Models.User.GetByEmail(request.Email)
	if err != nil {
		app.ErrorJSON(w, errors.New("Invalid credentials"), http.StatusUnauthorized)
		return
	}

	valid, err := user.PasswordMatches(request.Password)
	if err != nil || !valid {
		app.ErrorJSON(w, errors.New("Invalid credentials"), http.StatusUnauthorized)
		return
	}

	payload := jsonResponse{
		Error:   false,
		Message: fmt.Sprintf("Authenticated as %s", user.Email),
		Data:    user,
	}

	app.WriteJSON(w, http.StatusAccepted, payload)
}

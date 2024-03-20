package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
)

type Request struct {
	Action      string      `json:"action"`
	AuthPayload AuthPayload `json:"auth_payload"`
}

type AuthPayload struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (app *Config) Broker(w http.ResponseWriter, r *http.Request) {
	payload := jsonResponse{
		Error:   false,
		Message: "Hit broker service",
	}

	app.WriteJSON(w, http.StatusOK, payload)
}

func (app *Config) HandleSubmission(w http.ResponseWriter, r *http.Request) {
	var request Request

	err := app.ReadJSON(w, r, &request)
	if err != nil {
		app.ErrorJSON(w, err)
		return
	}

	switch request.Action {
	case "auth":
		app.Authenticate(w, request.AuthPayload)
	default:
		app.ErrorJSON(w, errors.New("invalid action"))
	}
}

func (app *Config) Authenticate(w http.ResponseWriter, a AuthPayload) {
	data, _ := json.MarshalIndent(a, "", "\t")

	request, err := http.NewRequest("POST", "http://authentication-broker/authenticate", bytes.NewBuffer(data))
	if err != nil {
		app.ErrorJSON(w, err)
		return
	}

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		app.ErrorJSON(w, err)
		return
	}
	defer response.Body.Close()

	if response.StatusCode == http.StatusUnauthorized {
		app.ErrorJSON(w, errors.New("Invalid credentials"))
		return
	}
	if response.StatusCode != http.StatusAccepted {
		app.ErrorJSON(w, errors.New("Error occurred"))
		return
	}

	var authResponse jsonResponse

	err = json.NewDecoder(response.Body).Decode(&authResponse)
	if err != nil {
		app.ErrorJSON(w, err)
		return
	}

	if authResponse.Error {
		app.ErrorJSON(w, err, http.StatusUnauthorized)
		return
	}

	var payload jsonResponse

	payload.Error = false
	payload.Message = "Authenticated"
	payload.Data = authResponse.Data

	app.WriteJSON(w, http.StatusAccepted, payload)
}

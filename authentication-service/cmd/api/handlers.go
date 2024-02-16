package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"
)

func (app *Config) Authenticate(w http.ResponseWriter, r *http.Request) {
	// Declare a variable that has same tags as in the json
	log.Println("DEBUG: Inside authenticate-service Authenticate")
	var requestPayload struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	err := app.readJSON(w, r, &requestPayload)

	if err != nil {
		app.errorJSON(w, err, http.StatusBadRequest)
		return
	}

	// Validate the user against the database
	user, err := app.Models.User.GetByEmail(requestPayload.Email)
	if err != nil {
		log.Println("DEBUG: Invalid Password")
		app.errorJSON(w, errors.New("invalid email credentials"), http.StatusBadRequest)
		return
	}

	valid, err := user.PasswordMatches(requestPayload.Password)
	if err != nil || !valid {
		log.Println("DEBUG: Invalid Password")
		app.errorJSON(w, errors.New("invalid password credentials"), http.StatusBadRequest)
		return
	}

	payload := jsonResponse{
		Error:   false,
		Message: fmt.Sprintf("Logged in user %s", user.Email),
		Data:    user,
	}

	app.writeJSON(w, http.StatusAccepted, payload)

}

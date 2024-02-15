package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
)

type RequestPayload struct {
	Action string      `json:"action"`
	Auth   AuthPayload `json:"auth,omitempty"`
}

type AuthPayload struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (app *Config) Broker(w http.ResponseWriter, r *http.Request) {

	payload := jsonResponse{
		Error:   false,
		Message: "Hit the broker",
	}

	_ = app.writeJSON(w, http.StatusOK, payload)
	// // Make it beautiful
	// out, _ := json.MarshalIndent(payload, "", "\t")

	// w.Header().Set("Content-Type", "application/json")
	// w.WriteHeader(http.StatusAccepted)
	// w.Write(out)

}

func (app *Config) HandleSubmission(w http.ResponseWriter, r *http.Request) {
	var requestPayload RequestPayload

	err := app.readJSON(w, r, &requestPayload) // read the input in requestPayload
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	switch requestPayload.Action {
	case "auth":
		// log.Println(fmt.Sprintf("%s"), r.Body)

		app.authenticate(w, requestPayload.Auth) // send the auth part of the payload (email,password)

	default:
		app.errorJSON(w, errors.New("unkown action"))
	}
}

func (app *Config) authenticate(w http.ResponseWriter, a AuthPayload) {
	// create some json we'll send to the auth microservice

	jsonData, _ := json.MarshalIndent(a, "", "\t") //format json

	// Create a new HTTP request to call authenticate service
	request, err := http.NewRequest("POST", "http://authentication-service/authenticate", bytes.NewBuffer(jsonData))
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	// Perform the HTTP request
	client := &http.Client{} // Create a HTTP client to perform the request

	response, err := client.Do(request) // Do sends the request and get the response

	if err != nil {
		app.errorJSON(w, err)
		return
	}

	defer response.Body.Close() //The response.Body represents the response body returned by the server. If you don't close the response body explicitly, it can lead to resource leaks, such as keeping the network connection open or consuming system resources unnecessarily.

	// make sure we get back the correct status code
	if response.StatusCode == http.StatusUnauthorized {
		app.errorJSON(w, errors.New("invalid credentials"))
		return
	} else if response.StatusCode != http.StatusAccepted {
		app.errorJSON(w, errors.New("Error calling auth service"))
		return
	}

	// Create a variable we'll read response.Body into
	var jsonFromService jsonResponse

	//Decode the json from the auth service and put in jsonfromService
	err = json.NewDecoder(response.Body).Decode(&jsonFromService)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	if jsonFromService.Error {
		app.errorJSON(w, err, http.StatusUnauthorized)
		return
	}

	// Send back to end user, data present in jsonFromService, which has data from 'response'

	var payload jsonResponse

	payload.Error = false
	payload.Message = "Authenticated"
	payload.Data = jsonFromService.Data

	app.writeJSON(w, http.StatusAccepted, payload) // Respond to the calling

}

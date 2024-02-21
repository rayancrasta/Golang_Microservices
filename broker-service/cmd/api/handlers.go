package main

import (
	"broker/event"
	"bytes"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"net/rpc"
)

type RequestPayload struct {
	Action string      `json:"action"`
	Auth   AuthPayload `json:"auth,omitempty"`
	Log    LogPayload  `json:"log,omitempty"`
	Mail   MailPayload `json:"mail,omitempty"`
}

type AuthPayload struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LogPayload struct {
	Name string `json:"name`
	Data string `json:"data"`
}

type MailPayload struct {
	FromAddress string `json:"from"`
	To          string `json:"to"`
	Subject     string `json:"subject"`
	Message     string `json:"message"`
}

func (app *Config) Broker(w http.ResponseWriter, r *http.Request) {

	log.Println("DEBUG: Inside Routes Broker function")

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
	log.Println("DEBUG: In Broker-HandleSubmission Function")

	var requestPayload RequestPayload

	err := app.readJSON(w, r, &requestPayload) // read the input in requestPayload
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	switch requestPayload.Action {
	case "auth":
		// log.Println(fmt.Sprintf("%s"), r.Body)
		log.Println("DEBUG: Switch case auth")
		app.authenticate(w, requestPayload.Auth) // send the auth part of the payload (email,password)
	case "log":
		log.Println("DEBUG: Switch case log")
		// app.logEventViaRabbit(w, requestPayload.Log)
		app.logItemviaRPC(w, requestPayload.Log)

	case "mail":
		log.Println("DEBUG: Switch case Mail")
		app.SendMail(w, requestPayload.Mail)
	default:
		app.errorJSON(w, errors.New("unkown action"))
	}
}

func (app *Config) logItem(w http.ResponseWriter, entry LogPayload) {
	log.Println("Inside Broker-logItem function")
	jsonData, _ := json.MarshalIndent(entry, "", "\t")
	logServiceURL := "http://logger-service/log"

	request, err := http.NewRequest("POST", logServiceURL, bytes.NewBuffer(jsonData))
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	request.Header.Set("Content-Type", "application/json")
	client := &http.Client{}

	response, err := client.Do(request)

	if err != nil {
		app.errorJSON(w, err)
		return
	}

	defer response.Body.Close()

	if response.StatusCode != http.StatusAccepted {
		app.errorJSON(w, err)
		return
	}

	var payload jsonResponse
	payload.Error = false
	payload.Message = "Logged"

	app.writeJSON(w, http.StatusAccepted, payload)
}

func (app *Config) authenticate(w http.ResponseWriter, a AuthPayload) {
	// create some json we'll send to the auth microservice
	log.Println("Inside Authenticate function")
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

	log.Println(response.StatusCode)
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

func (app *Config) SendMail(w http.ResponseWriter, msg MailPayload) {

	jsonData, _ := json.MarshalIndent(msg, "", "\t")

	log.Println("Broker mailer json data:", bytes.NewBuffer(jsonData))
	// Call the mail service
	mailServiceURL := "http://mailer-service/send"

	// Post to mail service
	request, err := http.NewRequest("POST", mailServiceURL, bytes.NewBuffer(jsonData))

	if err != nil {
		app.errorJSON(w, err)
		return
	}

	request.Header.Set("Content-Type", "application/json")

	client := &http.Client{}

	response, err := client.Do(request)

	if err != nil {
		app.errorJSON(w, err)
		return
	}

	defer response.Body.Close()

	// Make sure we get the right status code

	if response.StatusCode != http.StatusAccepted {
		app.errorJSON(w, errors.New("error calling mail service"))
		return
	}

	// Send back json
	var payload jsonResponse
	payload.Error = false
	payload.Message = "Message sent to " + msg.To

	app.writeJSON(w, http.StatusAccepted, payload)
}

func (app *Config) logEventViaRabbit(w http.ResponseWriter, logpayload LogPayload) {
	err := app.pushToQueue(logpayload.Name, logpayload.Data)

	if err != nil {
		app.errorJSON(w, err)
		return
	}

	var payload jsonResponse
	payload.Error = false
	payload.Message = "Logged via RabbitMQ"

	app.writeJSON(w, http.StatusAccepted, payload) // return to webUI

}

func (app *Config) pushToQueue(name, msg string) error {
	emitter, err := event.NewEventEmitter(app.Rabbit) // get amqp connection
	if err != nil {
		return err
	}

	payload := LogPayload{
		Name: name,
		Data: msg,
	}

	jsonData, _ := json.MarshalIndent(&payload, "", "\t") // get in JSON format
	err = emitter.Push(string(jsonData), "log.INFO")      // log.INFO is severity, routing key

	if err != nil {
		return err
	}

	return nil
}

type RPCPayload struct {
	Name string
	Data string
}

func (app *Config) logItemviaRPC(w http.ResponseWriter, logpayload LogPayload) {
	client, err := rpc.Dial("tcp", "logger-service:5001") // target RPC server
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	// Create a type that exactly matches, the remote RPC server is expecting to get
	rpcPayload := RPCPayload{
		Name: logpayload.Name,
		Data: logpayload.Data,
	}

	var result string
	err = client.Call("RPCServer.LogInfo", // Same as defined in remote rPC server
		rpcPayload,
		&result)

	if err != nil {
		app.errorJSON(w, err)
		return
	}

	var payload jsonResponse
	payload.Error = false
	payload.Message = result

	app.writeJSON(w, http.StatusAccepted, payload) // return to webUI
}

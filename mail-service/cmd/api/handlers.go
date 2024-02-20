package main

import (
	"log"
	"net/http"
)

func (app *Config) SendMail(w http.ResponseWriter, r *http.Request) {

	type mailMessage struct {
		From    string `json:"from"`
		To      string `json:"to"`
		Subject string `json:"subject"`
		Message string `json:"message"`
	}

	var requestPayload mailMessage

	err := app.readJSON(w, r, &requestPayload)

	if err != nil {
		log.Println("DEBUG: ReadJSON MailServer error: ", err)
		app.errorJSON(w, err)
	}

	log.Println("======Request payload MAIL DETAILS=========START")
	log.Println("From: ", requestPayload.From)
	log.Println("To: ", requestPayload.To)
	log.Println("Subject: ", requestPayload.Subject)
	log.Println("Plain Message", requestPayload.Message)
	log.Println("MESSAGE DETAILS++++++++ END")

	msg := Message{
		From:    requestPayload.From,
		To:      requestPayload.To,
		Subject: requestPayload.Subject,
		Data:    requestPayload.Message,
	}

	err = app.Mailer.SendSMTPMessage(msg)
	if err != nil {
		log.Println("DEBUG: SendSMTP error: ", err)
		app.errorJSON(w, err)
	}

	payload := jsonResponse{
		Error:   false,
		Message: "Sent to " + requestPayload.To,
	}

	app.writeJSON(w, http.StatusAccepted, payload)

}

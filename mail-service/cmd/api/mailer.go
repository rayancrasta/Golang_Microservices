package main

import (
	"bytes"
	"html/template"
	"log"
	"time"

	"github.com/vanng822/go-premailer/premailer"
	mail "github.com/xhit/go-simple-mail/v2"
)

type Mail struct { // Mail server config
	Domain      string
	Host        string
	Port        int
	Username    string
	Password    string
	Encryption  string
	FromAddress string
	FromName    string
}

type Message struct {
	From        string
	FromName    string
	To          string
	Subject     string
	Attachments []string
	Data        any
	DataMap     map[string]any
}

// Function to send email

func (m *Mail) SendSMTPMessage(msg Message) error {
	if msg.From == "" {
		msg.From = m.FromAddress
	}

	if msg.FromName == "" {
		msg.FromName = m.FromName
	}

	data := map[string]any{
		"message": msg.Data,
	}

	msg.DataMap = data

	formattedMessage, err := m.buildHTMLMessage(msg)
	if err != nil {
		return err
	}
	// HTML to plaintext

	plainMessage, err := m.buildPlainTextMessage(msg)
	if err != nil {
		return err
	}

	//Mail server
	server := mail.NewSMTPClient()
	server.Host = m.Host
	server.Port = m.Port
	server.Username = m.Username
	server.Password = m.Password
	server.Encryption = m.getEncryption(m.Encryption) // for production
	server.KeepAlive = false
	server.ConnectTimeout = 10 * time.Second
	server.SendTimeout = 10 * time.Second

	smtpClient, err := server.Connect()

	if err != nil {
		log.Println("DEBUG: Mail server connection error : ", err)
		return err
	}

	email := mail.NewMSG()
	email.SetFrom(msg.From).
		AddTo(msg.To).
		SetSubject(msg.Subject)

	email.SetBody(mail.TextPlain, plainMessage)
	email.AddAlternative(mail.TextHTML, formattedMessage)

	if len(msg.Attachments) > 0 {
		for _, x := range msg.Attachments {
			email.AddAttachment(x)
		}
	}

	log.Println("======MAIL DETAILS=========START")
	log.Println("From: ", msg.From)
	log.Println("To: ", msg.To)
	log.Println("Subject: ", msg.Subject)
	log.Println("Plain Message", plainMessage)
	log.Println("HTML message", formattedMessage)
	log.Println("MESSAGE DETAILS++++++++ END")

	err = email.Send(smtpClient)
	if err != nil {
		log.Println("DEBUG: SendSMTP Mailer.go error: ", err)
		return err
	}

	return nil
}

func (m *Mail) buildPlainTextMessage(msg Message) (string, error) {

	templatetoRender := "./templates/mail.plain.gohtml"

	t, err := template.New("email-plain").ParseFiles(templatetoRender)
	if err != nil {
		return "", err
	}

	var tpl bytes.Buffer

	//Set our email data into the template
	if err = t.ExecuteTemplate(&tpl, "body", msg.DataMap); err != nil {
		return "", err
	}

	plainMessage := tpl.String()

	return plainMessage, nil

}

func (m *Mail) buildHTMLMessage(msg Message) (string, error) {

	templatetoRender := "./templates/mail.html.gohtml"

	t, err := template.New("email.html").ParseFiles(templatetoRender)
	if err != nil {
		return "", err
	}

	var tpl bytes.Buffer

	//Set our email data into the template
	if err = t.ExecuteTemplate(&tpl, "body", msg.DataMap); err != nil {
		return "", err
	}

	formattedMessage := tpl.String()

	//Inline the CSS
	// formattedMessage is html
	formattedMessage, err = m.inlineCSS(formattedMessage)
	if err != nil {
		return "", err
	}

	return formattedMessage, nil

}

func (m *Mail) inlineCSS(s string) (string, error) {

	// Use the 3rd party premailer package
	options := premailer.Options{
		RemoveClasses:     false,
		CssToAttributes:   false,
		KeepBangImportant: true,
	}

	prem, err := premailer.NewPremailerFromString(s, &options)

	if err != nil {
		return "", err
	}

	html, err := prem.Transform()

	if err != nil {
		return "", err
	}

	return html, nil
}

func (m *Mail) getEncryption(s string) mail.Encryption {
	switch s {
	case "tls":
		return mail.EncryptionSTARTTLS
	case "ssl":
		return mail.EncryptionSSLTLS
	case "none", "":
		return mail.EncryptionNone
	default:
		return mail.EncryptionSTARTTLS
	}
}

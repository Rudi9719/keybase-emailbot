package main

import "fmt"
import "net/smtp"

func send(e Email) {
	e.Body = signMessage(e.Body)

	message := fmt.Sprintf("From: %s\n", conf.MyEmail)
	for _, recipient := range e.Recipients {
		message += fmt.Sprintf("To: %s\n", recipient)
	}
	for _, cc := range e.Cc {
		message += fmt.Sprintf("Cc: %s\n", cc)
	}
	for _, bcc := range e.Bcc {
		message += fmt.Sprintf("Bcc: %s\n", bcc)
	}
	message += fmt.Sprintf("Subject: %s\n", e.Subject)
	message += e.Body
	log.LogInfo("Message created")
	log.LogDebug(message)
	log.LogInfo("Sending message")
	err := smtp.SendMail(conf.SmtpServer,
		smtp.PlainAuth("", conf.MyEmail, conf.EmailPass, conf.AuthServer),
		conf.MyEmail, e.Recipients, []byte(message))
	if err != nil {
		log.LogErrorType(err)
	}
	log.LogInfo("Email Sent")
}

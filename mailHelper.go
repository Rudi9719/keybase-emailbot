package main

import "fmt"
import "io/ioutil"
import "net/smtp"
import "samhofi.us/x/keybase"

func send(e Email, api keybase.ChatAPI) {
	e.Body = appendSignature(e.Body)

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
	go chat.React(api.Msg.ID, ":mailbox_with_no_mail:")
	if conf.KeyPass == "" {
		go chat.React(api.Msg.ID, ":unlock:")
	} else {
		go chat.React(api.Msg.ID, ":lock_with_ink_pen:")
	}
	err := smtp.SendMail(conf.SmtpServer,
		smtp.PlainAuth("", conf.MyEmail, conf.EmailPass, conf.AuthServer),
		conf.MyEmail, e.Recipients, []byte(message))
	go chat.React(api.Msg.ID, ":mailbox_with_no_mail:")
	if err != nil {
		log.LogErrorType(err)
		chat.React(api.Msg.ID, ":warning:")
		return
	}
	go chat.React(api.Msg.ID, ":mailbox_with_mail:")
	log.LogInfo("Email Sent")
}

func appendSignature(body string) string {
	bytes, err := ioutil.ReadFile("default.sig")
	if err != nil {
		log.LogErrorType(err)
	}
	return signMessage(fmt.Sprintf("%s\n\n%s\n", body, string(bytes)))

}

package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strings"
	"syscall"

	"github.com/rudi9719/loggy"
	"golang.org/x/crypto/ssh/terminal"
	"samhofi.us/x/keybase"
)

var (
	k = keybase.NewKeybase()

	logOpts = loggy.LogOpts{
		//OutFile:   "Keybase-Email.log",
		//KBTeam:    "nightmarehaus.logs",
		//KBChann:   "general",
		//ProgName:  "KB-Email",
		Level:     5,
		UseStdout: true,
	}

	log  = loggy.NewLogger(logOpts)
	conf = Config{}
)

func main() {
	if !k.LoggedIn {
		log.LogPanic("Keybase not logged in.")
	}
	log.LogInfo(fmt.Sprintf("Bot started using account %s", k.Username))
	conf = loadConfig()
	setupCredentials()
	log.LogInfo("Starting keybase")
	k.Run(func(api keybase.ChatAPI) {
		handleMessage(api)
	})

}

func handleMessage(api keybase.ChatAPI) {
	if api.Msg.Channel.Name != k.Username {
		log.LogInfo("Wrong channel detected.")
		return
	}
	if api.Msg.Sender.Username != k.Username {
		log.LogInfo("Wrong username detected.")
		return
	}
	if api.Msg.Content.Type != "text" {
		log.LogInfo("Wrong message type detected.")
		return
	}
	parts := strings.Split(api.Msg.Content.Text.Body, " ")
	if parts[0] != "!email" {
		log.LogInfo("Wrong command detected")
		return
	}
	if len(parts) < 4 {
		log.LogInfo("Wrong length of parts detected.")
		return
	}
	var e Email
	partCounter := 1
	for _, subj := range parts[1:] {
		if strings.Contains(subj, "@") {
			break
		}
		partCounter++
		e.Subject += fmt.Sprintf("%s ", subj)
	}
	for _, to := range parts {
		if strings.HasPrefix(to, "to:") {
			if strings.Contains(to, "@") {
				e.Recipients = append(e.Recipients, strings.Replace(to, "to:", "", -1))
				partCounter++
			}
		}
	}
	for _, cc := range parts {
		if strings.HasPrefix(cc, "cc:") {
			if strings.Contains(cc, "@") {
				e.Cc = append(e.Cc, strings.Replace(cc, "cc:", "", -1))
				partCounter++
			}
		}
	}
	for _, bcc := range parts {
		if strings.HasPrefix(bcc, "bcc:") {
			if strings.Contains(bcc, "@") {
				e.Bcc = append(e.Bcc, strings.Replace(bcc, "bcc:", "", -1))
				partCounter++
			}
		}
	}
	for _, word := range parts[partCounter:] {
		e.Body += fmt.Sprintf("%s ", word)
	}
	log.LogDebug(fmt.Sprintf("%+v", e))
	go send(e)
}

func loadConfig() Config {
	var c Config
	bytes, err := ioutil.ReadFile("conf.json")
	if err != nil {
		log.LogErrorType(err)
	}
	err = json.Unmarshal(bytes, &c)
	if err != nil {
		log.LogErrorType(err)
	}
	bytes, err = ioutil.ReadFile("priv.key")
	if err != nil {
		log.LogErrorType(err)
	}
	c.PrivateKey = string(bytes)
	return c
}
func setupCredentials() {
	log.LogCritical("Enter pgp key passphrase:")
	bytePass, err := terminal.ReadPassword(int(syscall.Stdin))
	if err != nil {
		log.LogCritical(fmt.Sprintf("Error reading pgp password:\n```%+v```", err))
	}
	conf.KeyPass = strings.TrimSpace(string(bytePass))
	log.LogCritical("Enter email passphrase:")
	bytePass, err = terminal.ReadPassword(int(syscall.Stdin))
	if err != nil {
		log.LogCritical(fmt.Sprintf("Error reading email password:\n```%+v```", err))
	}
	conf.EmailPass = strings.TrimSpace(string(bytePass))
}

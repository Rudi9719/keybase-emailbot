package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strings"
	"syscall"
	"time"

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
	em := Email{
		Recipients: []string{"rudi@nmare.net"},
		Subject:    "Test Email",
		Body:       "Hello, world!",
	}
	send(em)
	time.Sleep(2 * time.Second)

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

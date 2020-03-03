package main

import (
	"bytes"
	//	"encoding/base64"
	"fmt"
	"golang.org/x/crypto/openpgp"
	"golang.org/x/crypto/openpgp/clearsign"
	"strings"
)

func getPrivateKey() *openpgp.Entity {
	pp := conf.KeyPass
	ppb := []byte(pp)
	log.LogDebug("Getting entityList")
	entitylist, err := openpgp.ReadArmoredKeyRing(strings.NewReader(conf.PrivateKey))
	if err != nil {
		log.LogErrorType(err)
	}
	log.LogDebug(fmt.Sprintf("Getting entity 0 ```%+v```", entitylist))
	entity := entitylist[0]
	log.LogDebug("if PrivateKey != nil")
	if entity.PrivateKey != nil && entity.PrivateKey.Encrypted {
		err := entity.PrivateKey.Decrypt(ppb)
		if err != nil {
			log.LogErrorType(err)
		}
	}

	for _, subkey := range entity.Subkeys {
		if subkey.PrivateKey != nil && subkey.PrivateKey.Encrypted {
			err := subkey.PrivateKey.Decrypt(ppb)
			if err != nil {
				log.LogErrorType(err)
			}
		}
	}
	return entity
}

func signMessage(m string) string {
	if conf.KeyPass == "" {
		return m
	}
	pk := getPrivateKey()
	out := new(bytes.Buffer)
	in, err := clearsign.Encode(out, pk.PrivateKey, nil)
	if err != nil {
		log.LogErrorType(err)
	}
	in.Write([]byte(m))
	in.Close()
	return out.String()
}

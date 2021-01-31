package main

import (
	"bytes"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/slack-go/slack"
)

func verifySigningSecret(r *http.Request) error {
	signingSecret := os.Getenv("SLACK_SIGNING_SECRET")
	verifier, err := slack.NewSecretsVerifier(r.Header, signingSecret)
	if err != nil {
		log.Println(err.Error())
		return err
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println(err.Error())
		return err
	}
	// Need to use r.Body again when unmarshalling SlashCommand and InteractionCallback
	r.Body = ioutil.NopCloser(bytes.NewBuffer(body))

	verifier.Write(body)
	if err = verifier.Ensure(); err != nil {
		log.Println(err.Error())
		return err
	}

	return nil
}

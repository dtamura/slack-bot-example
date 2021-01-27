package main

import (
	"io"
	"io/ioutil"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	opentracing "github.com/opentracing/opentracing-go"
	"github.com/slack-go/slack"
)

func slash(c *gin.Context) {
	span, _ := opentracing.StartSpanFromContext(c.Request.Context(), "slash")

	defer span.Finish()
	verifier, err := slack.NewSecretsVerifier(c.Request.Header, os.Getenv("SLACK_SIGNING_SECRET"))
	if err != nil {
		c.Error(err)
		c.AbortWithStatusJSON(401, gin.H{"status": false, "message": err.Error()})
		return
	}

	c.Request.Body = ioutil.NopCloser(io.TeeReader(c.Request.Body, &verifier))
	sc, err := slack.SlashCommandParse(c.Request)
	if err != nil {
		c.Error(err)
		c.AbortWithStatusJSON(401, gin.H{"status": false, "message": err.Error()})
		return
	}

	if err = verifier.Ensure(); err != nil {
		c.Error(err)
		c.AbortWithStatusJSON(401, gin.H{"status": false, "message": err.Error()})
		return
	}

	obj := slack.NewTextBlockObject(slack.MarkdownType, "<@"+sc.UserName+"> Hello World!!", false, false)
	block := slack.NewSectionBlock(obj, nil, nil)

	modal := slack.ModalViewRequest{
		Type:   slack.VTModal,
		Title:  slack.NewTextBlockObject(slack.PlainTextType, "hoge", true, false),
		Close:  slack.NewTextBlockObject(slack.PlainTextType, "close", true, false),
		Submit: slack.NewTextBlockObject(slack.PlainTextType, "submit", true, false),
		Blocks: slack.Blocks{BlockSet: []slack.Block{block}},
	}
	log.Println(modal)

	hostname, _ := os.Hostname()

	span.SetTag("hostname", hostname)

	// 成功時
	_, err = slackAPI.OpenView(sc.TriggerID, modal)
	if err != nil {
		c.Error(err)
		c.AbortWithStatusJSON(500, gin.H{"status": false, "message": err.Error()})
	}
	return
	// c.String(200, "<@"+sc.UserName+"> Hello World!!")
}

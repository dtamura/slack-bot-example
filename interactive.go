package main

import (
	"encoding/json"
	"log"

	"github.com/gin-gonic/gin"
	opentracing "github.com/opentracing/opentracing-go"
	"github.com/slack-go/slack"
)

func interactive(c *gin.Context) {
	span, _ := opentracing.StartSpanFromContext(c.Request.Context(), "interactive")
	defer span.Finish()

	if err := verifySigningSecret(c.Request); err != nil {
		c.Error(err)
		c.AbortWithStatusJSON(401, gin.H{"status": false})
		return
	}

	var ic slack.InteractionCallback
	err := json.Unmarshal([]byte(c.Request.FormValue("payload")), &ic)
	if err != nil {
		log.Printf("Could not parse action response JSON: %v", err)
		c.Error(err)
		c.AbortWithStatusJSON(500, gin.H{"status": false, "message": err.Error()})
		return
	}
	switch ic.Type {
	case slack.InteractionTypeShortcut:
		handleShortcut(c, ic)
	case slack.InteractionTypeBlockActions:
		handleBlockAction(c, ic)
	case slack.InteractionTypeViewSubmission:
		handleViewSubmission(c, ic)
	}

	return
}

func handleViewSubmission(c *gin.Context, ic slack.InteractionCallback) error {
	log.Printf("view id: " + ic.View.ID)

	// msg := fmt.Sprintf("Hello %s, nice to meet you!", ic.User.ID)

	res := slack.NewClearViewSubmissionResponse()
	c.JSON(200, res)

	// _, _, err := slackAPI.PostMessage(ic.User.ID,
	// 	slack.MsgOptionText(msg, false),
	// 	slack.MsgOptionAttachments())
	// if err != nil {
	// 	c.Error(err)
	// 	c.AbortWithStatusJSON(401, gin.H{"status": false})
	// 	return err
	// }
	return nil
}

func handleBlockAction(c *gin.Context, ic slack.InteractionCallback) error {
	log.Printf("AcrtionID: %v, Value: %v", ic.ActionCallback.BlockActions[0].Type, ic.ActionCallback.BlockActions[0].Value)
	by, _ := ic.MarshalJSON()
	log.Printf(string(by))
	for _, v := range ic.ActionCallback.BlockActions {
		switch v.Value {
		case "click_me_123":
			updateView(c, ic)
		default:
		}
	}
	c.String(200, "<@"+ic.User.ID+"> Hello World!!")
	return nil
}

func updateView(c *gin.Context, ic slack.InteractionCallback) error {
	obj := slack.NewTextBlockObject(slack.MarkdownType, "*This* is a section block with a button.", false, false)
	block := slack.NewSectionBlock(obj,
		nil,
		slack.NewAccessory(
			slack.NewButtonBlockElement("click_me", "click_me_123",
				slack.NewTextBlockObject(slack.PlainTextType, "Click Me", true, false))),
	)
	block2 := slack.NewActionBlock(
		"block2",
		slack.NewOptionsSelectBlockElement(
			slack.OptTypeChannels,
			slack.NewTextBlockObject(
				slack.PlainTextType,
				"select channel",
				true,
				false,
			),
			"select_channel",
		),
	)
	modal := slack.ModalViewRequest{
		Type:   slack.VTModal,
		Title:  slack.NewTextBlockObject(slack.PlainTextType, "hoge", true, false),
		Close:  slack.NewTextBlockObject(slack.PlainTextType, "close", true, false),
		Submit: slack.NewTextBlockObject(slack.PlainTextType, "submit", true, false),
		Blocks: slack.Blocks{BlockSet: []slack.Block{
			block,
			block2,
		}},
	}
	byt, _ := json.Marshal(modal)
	log.Println(string(byt))

	resp, err := slackAPI.UpdateView(modal, ic.View.ExternalID, ic.View.Hash, ic.Container.ViewID)
	if err != nil {
		log.Panic(err)
	}
	byt, _ = json.Marshal(resp)
	log.Println(byt)

	return nil
}

func handleShortcut(c *gin.Context, ic slack.InteractionCallback) error {
	log.Printf("callback_id: " + ic.CallbackID)
	obj := slack.NewTextBlockObject(slack.MarkdownType, "*This* is a section block with a button.", false, false)
	block := slack.NewSectionBlock(obj,
		nil,
		slack.NewAccessory(
			slack.NewButtonBlockElement("click_me", "click_me_123",
				slack.NewTextBlockObject(slack.PlainTextType, "Click Me", true, false))),
	)
	block2 := slack.NewActionBlock(
		"block2",
		slack.NewOptionsSelectBlockElement(
			slack.OptTypeChannels,
			slack.NewTextBlockObject(
				slack.PlainTextType,
				"select channel",
				true,
				false,
			),
			"select_channel",
		),
	)
	block3 := slack.NewActionBlock(
		"block3",
		slack.NewOptionsSelectBlockElement(
			slack.OptTypeStatic,
			slack.NewTextBlockObject(
				slack.PlainTextType,
				"select hoge",
				true,
				false,
			),
			"select_hoge",
			slack.NewOptionBlockObject(
				"hoge",
				slack.NewTextBlockObject(slack.PlainTextType, "hoge", false, false),
				slack.NewTextBlockObject(slack.PlainTextType, "hoge", false, false),
			),
		),
	)

	modal := slack.ModalViewRequest{
		Type:   slack.VTModal,
		Title:  slack.NewTextBlockObject(slack.PlainTextType, "hoge", true, false),
		Close:  slack.NewTextBlockObject(slack.PlainTextType, "close", true, false),
		Submit: slack.NewTextBlockObject(slack.PlainTextType, "submit", true, false),
		Blocks: slack.Blocks{BlockSet: []slack.Block{
			block,
			block2,
			block3,
		}},
	}
	byt, _ := json.Marshal(modal)
	log.Println(string(byt))

	_, err := slackAPI.OpenView(ic.TriggerID, modal)

	if err != nil {
		log.Panic(err)
		return err
	}
	return nil
}

package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"net/url"
	
	ce "github.com/cloudevents/sdk-go"
	"github.com/bwmarrin/discordgo"
)


type eventReceiver struct{}

func (r *eventReceiver) Receive(ctx context.Context, event ce.Event, resp *ce.EventResponse) error {

	log.Printf("Raw Event: %v", event)

	if event.ID() == "" {
		log.Println("unable to parse event ID")
		return errors.New("invalid event format")
	}

	var mdata map[string]interface{}
	if err := event.DataAs(&mdata); err != nil {
		log.Printf("Failed to DataAs: %s", err.Error())
		return err
	}

	var webhook *discordgo.Webhook = mdata["webhook"].(*discordgo.Webhook)

	var posturl = "https://discordapp.com/webhooks/" + webhook.ID +"/" + webhook.Token

	var data = url.Values{
		"content" : {"pong!"},
	}

	_,err := http.PostForm(posturl,data)
	if err != nil {
		log.Printf("webhook post failed: %s", err.Error())
	}

	re := &ce.EventResponse{
		Status:  200,
		Event:   &event,
		Reason:  "Done!!!",
		Context: event.Context,
	}

	resp = re

	return nil

}

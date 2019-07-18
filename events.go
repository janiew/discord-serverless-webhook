package main

import (
	"context"
	"errors"
	"log"
	"strings"

	"github.com/PratikMahajan/Twitter-Serverless-Serving/config"
	ce "github.com/cloudevents/sdk-go"
	"github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/oauth1"
)

var (
	consumerKey    = strings.TrimSpace(config.MustGetEnvVar("T_CONSUMER_KEY", ""))
	consumerSecret = strings.TrimSpace(config.MustGetEnvVar("T_CONSUMER_SECRET", ""))
	accessToken    = strings.TrimSpace(config.MustGetEnvVar("T_ACCESS_TOKEN", ""))
	accessSecret   = strings.TrimSpace(config.MustGetEnvVar("T_ACCESS_SECRET", ""))
)


type eventReceiver struct{}

func (r *eventReceiver) Receive(ctx context.Context, event ce.Event, resp *ce.EventResponse) error {

	// twitter client config
	cfg := oauth1.NewConfig(consumerKey, consumerSecret)
	token := oauth1.NewToken(accessToken, accessSecret)
	httpClient := cfg.Client(oauth1.NoContext, token)
	twClient := twitter.NewClient(httpClient)


	log.Printf("Raw Event: %v", event)

	if event.ID() == "" {
		log.Println("unable to parse event ID")
		return errors.New("invalid event format")
	}


	var tdata map[string]interface{}
	if err := event.DataAs(&tdata); err != nil {
		log.Printf("Failed to DataAs: %s", err.Error())
		return err
	}

	tweet, _, err := twClient.Statuses.Update("just setting up my twttr", nil)
	if err!=nil{
		log.Printf("failed to send Tweet : %s", err.Error())
	}
	log.Printf("sent tweet: %s\n", tweet.IDStr)

	re := &ce.EventResponse{
		Status:  200,
		Event:   &event,
		Reason:  "Done!!!",
		Context: event.Context,
	}

	resp = re

	return nil

}

package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/PratikMahajan/Twitter-Serverless-Serving/config"
	"github.com/PratikMahajan/Twitter-Serverless-Serving/logic"
	ce "github.com/cloudevents/sdk-go"
	"github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/oauth1"
)

var (
	consumerKey    = strings.TrimSpace(config.MustGetEnvVar("T_CONSUMER_KEY", ""))
	consumerSecret = strings.TrimSpace(config.MustGetEnvVar("T_CONSUMER_SECRET", ""))
	accessToken    = strings.TrimSpace(config.MustGetEnvVar("T_ACCESS_TOKEN", ""))
	accessSecret   = strings.TrimSpace(config.MustGetEnvVar("T_ACCESS_SECRET", ""))
	selfHandle     = strings.TrimSpace(config.MustGetEnvVar("SELF_HANDLE",""))
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



	// Get Tweet ID
	tID := fmt.Sprintf("%v", tdata["id_str"])

	tIDINT, err := strconv.ParseInt(tID, 10, 64)
	if err != nil {
		log.Printf("Unable to get tweet id: %s", err.Error())
	}


	// Get User Twitter Handle

	datajson, err := json.Marshal(tdata["user"])
	if err != nil {
		log.Printf("Unable to get user json: %s", err.Error())
	}

	type userJson struct{
		ScreenName string `json:"screen_name"`
	}


	// convert to json
	var userData userJson
	err =json.Unmarshal(datajson, &userData)
	if err != nil {
		log.Printf("json unmarshal failed: %s", err.Error())
	}

	userHandle := userData.ScreenName

	if selfHandle != strings.ToLower(userHandle) {

		// Setting up the Reply params
		params := &twitter.StatusUpdateParams{
			InReplyToStatusID: tIDINT,
		}

		// Get Video URL
		videoUrl := logic.Extract_url(ctx, tdata["quoted_status"])

		var tweetData string
		if videoUrl != "" {
			// Composing tweet Data
			tweetData = "Copy and Paste URL in Browser To Download "+ videoUrl +" " + "twitter.com/" + userHandle + "/status/" + string(tID)
		} else {
			tweetData = "No Video Found"
		}

		// Sending Tweet
		_, _, errStatus := twClient.Statuses.Update(tweetData, params)
		if errStatus != nil {
			log.Printf("failed to send Tweet : %s", errStatus.Error())
		}
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

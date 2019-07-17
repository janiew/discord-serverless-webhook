package main

import (
	"context"
	"errors"
	"log"

	ce "github.com/cloudevents/sdk-go"
)


type eventReceiver struct{}

func (r *eventReceiver) Receive(ctx context.Context, event ce.Event, resp *ce.EventResponse) error {

	log.Printf("Raw Event: %v", event)

	if event.ID() == "" {
		log.Println("unable to parse event ID")
		return errors.New("invalid event format")
	}


	var p map[string]interface{}
	if err := event.DataAs(&p); err != nil {
		log.Printf("Failed to DataAs: %s", err.Error())
		return err
	}

	ed := &eventData{Context: event.Context.AsV02(), Data: p}

	re := &ce.EventResponse{
		Status:  200,
		Event:   &event,
		Reason:  "Stored",
		Context: event.Context,
	}

	resp = re

	return nil

}

package logic

import (
	"context"
	"encoding/json"
	"log"
)

func Extract_url(ctx context.Context, data interface{}) string {

	// Get User Twitter Handle

	dataJson, err := json.Marshal(data)
	if err != nil {
		log.Printf("Unable to get video json: %s", err.Error())
		return ""
	}

	// convert to json
	var quoteS QuotedStatus
	err =json.Unmarshal(dataJson, &quoteS)
	if err != nil {
		log.Printf("json unmarshal failed for video entities: %s", err.Error())
		return ""
	}

	var videoEntities []Variant
	videoEntities = quoteS.ExtendedEntities.Media[0].VideoInfo.Variants

	videoUrl := ""
	maxBitrate := 0
	for _, video := range videoEntities {

		if video.Bitrate > maxBitrate{
			videoUrl = video.URL
		}

	}

	return videoUrl

}

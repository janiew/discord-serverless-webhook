package logic

type QuotedStatus struct {
		Coordinates        interface{} `json:"coordinates"`
		CreatedAt          string      `json:"created_at"`
		CurrentUserRetweet interface{} `json:"current_user_retweet"`
		ExtendedEntities   struct {
			Media []struct {
				VideoInfo struct {
					Variants []Variant`json:"variants"`
				} `json:"video_info"`
			} `json:"media"`
		} `json:"extended_entities"`
	}

type Variant struct {
	ContentType string `json:"content_type"`
	Bitrate     int    `json:"bitrate"`
	URL         string `json:"url"`
}
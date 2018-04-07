package model

import (
	"encoding/json"

	"github.com/max0ne/twitter_thing/back/db"
)

// Tweet - -
type Tweet struct {
	Tid       string `json:"tid"`
	UID       string `json:"uid"`
	Content   string `json:"content"`
	CreatedAt string `json:"createdAt"`
}

// NewTweet - -
func NewTweet(uid, content string) Tweet {
	return Tweet{
		UID:     uid,
		Content: content,
		// CreatedAt: time.Now(),
	}
}

// GetTweet - -
func GetTweet(tid string, table *db.Table) (*Tweet, error) {
	var tweet Tweet
	if err := json.Unmarshal([]byte(table.Get(tid)), &tweet); err != nil {
		return nil, err
	}
	return &tweet, nil
}

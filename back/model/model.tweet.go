package model

import (
	"encoding/json"
	"strings"

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

// getPostedBy return list of `tid`s posted by `vid`
func getPostedBy(vid string, postedByTable *db.Table) []string {
	return strings.Split(postedByTable.Get(vid), ",")
}

func getBucket(uid string, bucketTable *db.Table) []string {
	return strings.Split(bucketTable.Get(uid), ",")
}

// GetFollowers - -
func GetFollowers(vid string, followTable *db.Table) []string {
	return strings.Split(followTable.Get(vid), ",")
}

// PublishNewTweet - -
func PublishNewTweet(tweet Tweet, followTable, tweetTable, bucketTable, postedByTable *db.Table) error {

	tweet.Tid = tweetTable.IncID()
	tweetJSONBytes, err := json.Marshal(tweet)
	if err != nil {
		return err
	}

	// 0. store content
	tweetTable.Put(tweet.Tid, string(tweetJSONBytes))

	// 1. 发给自己的tweet里
	postedBy := postedByTable.Get(tweet.UID)
	postedBy += "," + tweet.Tid
	postedByTable.Put(tweet.UID, postedBy)

	// 2. 发给followers的buckets里
	followers := GetFollowers(tweet.UID, followTable)
	for _, follower := range followers {
		bucket := bucketTable.Get(follower)
		bucket += "," + tweet.Tid
		bucketTable.Put(follower, bucket)
	}
	return nil
}

package model

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/max0ne/twitter_thing/back/db"
)

// Tweet - -
type Tweet struct {
	Tid       string `json:"tid"`
	Uname     string `json:"uname"`
	Content   string `json:"content"`
	CreatedAt string `json:"createdAt"`
}

// NewTweet - -
func NewTweet(uname, content string) Tweet {
	return Tweet{
		Uname:   uname,
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

// DelTweet - -
func DelTweet(tid string, table *db.Table) error {
	table.Del(tid)
	return nil
}

// getPostedBy return list of `tid`s posted by `vid`
func getPostedBy(vid string, postedByTable *db.Table) []string {
	return strings.Split(postedByTable.Get(vid), ",")
}

type tweetBucket struct {
	uname string
	tid   string
}

func (tb tweetBucket) toString() string {
	return fmt.Sprintf("%s_%s", tb.uname, tb.tid)
}

func tweetBucketFromString(ss string) tweetBucket {
	var tb tweetBucket
	json.Unmarshal([]byte(ss), &tb)
	return tb
}

func getBucket(uname string, bucketTable *db.Table) ([]tweetBucket, error) {
	var tbs []tweetBucket
	err := bucketTable.GetObj(uname, &tbs)
	return tbs, err
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
	postedBy := postedByTable.Get(tweet.Uname)
	postedBy += "," + tweet.Tid
	postedByTable.Put(tweet.Uname, postedBy)

	// 2. 发给followers的buckets里
	followers := GetFollowers(tweet.Uname, followTable)
	newBucketItem := tweetBucket{tid: tweet.Tid, uname: tweet.Uname}
	for _, follower := range followers {
		buckets, err := getBucket(tweet.Uname, bucketTable)
		if err != nil {
			fmt.Println(err)
			continue
		}
		buckets = append(buckets, newBucketItem)
		bucketTable.PutObj(follower, buckets)
	}
	return nil
}

// UnfollowUserTweet remove tweets of `vid` from `user.uname`'s buckets
func UnfollowUserTweet(user User, vid string, bucketTable *db.Table) error {
	buckets, err := getBucket(user.Uname, bucketTable)
	if err != nil {
		return err
	}

	newBuckets := []tweetBucket{}
	for _, buck := range buckets {
		if buck.uname != vid {
			newBuckets = append(newBuckets, buck)
		}
	}
	return bucketTable.PutObj(user.Uname, newBuckets)
}

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
func DelTweet(tid string, tweetTable *db.Table) error {
	// only delete from tweet content in tweet table
	// don't delete from bucket bc there's no way to fast access
	// all buckets containing this tweet
	tweetTable.Del(tid)
	return nil
}

// getPostedBy return list of `tid`s posted by `vname`
func getPostedBy(vname string, postedByTable *db.Table) []string {
	return strings.Split(postedByTable.Get(vname), ",")
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

// PublishNewTweet - -
func PublishNewTweet(tweet Tweet, followerTable, tweetTable, bucketTable, postedByTable *db.Table) error {

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
	followers := GetFollowers(tweet.Uname, followerTable)
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

// UnfollowUserTweet remove tweets of `vname` from `user.uname`'s buckets
func UnfollowUserTweet(user User, vname string, bucketTable *db.Table) error {
	buckets, err := getBucket(user.Uname, bucketTable)
	if err != nil {
		return err
	}

	newBuckets := []tweetBucket{}
	for _, buck := range buckets {
		if buck.uname != vname {
			newBuckets = append(newBuckets, buck)
		}
	}
	return bucketTable.PutObj(user.Uname, newBuckets)
}

// GetUserTweets get all tweets from a specific user
func GetUserTweets(vname string, tweetTable, postedByTable *db.Table) []Tweet {
	tids := getPostedBy(vname, postedByTable)
	tweets := []Tweet{}
	for _, tid := range tids {
		tweet, err := GetTweet(tid, tweetTable)
		if tweet != nil && err == nil {
			tweets = append(tweets, *tweet)
		}
	}
	return tweets
}

// GetUserFeed get user's following users's tweets
func GetUserFeed(uname string, tweetTable, bucketTable *db.Table) ([]Tweet, error) {
	tbs, err := getBucket(uname, bucketTable)
	if err != nil {
		return []Tweet{}, err
	}

	tweets := []Tweet{}
	for _, tb := range tbs {
		tweet, err := GetTweet(tb.tid, tweetTable)
		if tweet == nil || err != nil {
			continue
		}
		tweets = append(tweets, *tweet)
	}
	return tweets, nil
}

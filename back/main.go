package main

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type User struct {
	Uid      int
	Username string
	Password string
}

type Tweet struct {
	Uid     int
	Tid     int
	Content string
	Time    string
}

var user_registered map[int]User

var t_tweet_content map[int]Tweet

var t_tweet_bucket map[int][]int

var t_posted_by map[int][]int

var t_follow map[int][]int

var usr_cnt int
var tweet_id int

func init() {
	usr_cnt = 0
	tweet_id = 0
	user_registered = make(map[int]User)
	t_tweet_content = make(map[int]Tweet)
	t_tweet_bucket = make(map[int][]int)
	t_posted_by = make(map[int][]int)
	t_follow = make(map[int][]int)
}

// RESTFul Apis
func signup(c *gin.Context) {
	username := c.PostForm("username")
	password := c.PostForm("password")

	// check whether the user already exists
	for _, u := range user_registered {
		if u.Username == username {
			c.JSON(http.StatusUnauthorized, gin.H{"status": "user already exists"})
			return
		}
	}
	usr_cnt += 1
	user := User{usr_cnt, username, password}
	user_registered[usr_cnt] = user

	c.JSON(200, gin.H{
		"status":   "posted",
		"username": username,
	})
}

func login(c *gin.Context) {
	username := c.PostForm("username")
	password := c.PostForm("password")

	for _, u := range user_registered {
		if u.Username == username && u.Password == password {
			c.JSON(200, gin.H{
				"status":   "posted",
				"username": username,
			})
			return
		}
		if u.Username == username {
			c.JSON(http.StatusUnauthorized, gin.H{"status": "wrong password"})
			return
		}
	}
	c.JSON(http.StatusUnauthorized, gin.H{"status": "user does not exist"})
}

func unregister(c *gin.Context) {
	// username := c.PostForm("username")
	// for _, u := range user_registered {
	// 	if u.Username == username && u.Password == password {
	// 		c.JSON(200, gin.H{
	// 			"status":   "posted",
	// 			"username": username,
	// 		})
	// 		return
	// 	}
	// 	if u.Username == username {
	// 		c.JSON(http.StatusUnauthorized, gin.H{"status": "wrong password"})
	// 		return
	// 	}
	// }
	// c.JSON(http.StatusUnauthorized, gin.H{"status": "user does not exist"})
}

func createNewTweet(c *gin.Context) {
	uid, err := strconv.Atoi(c.PostForm("uid"))
	if err != nil {
		fmt.Println(err)
	}
	content := c.PostForm("content")
	tweet_id = tweet_id + 1

	tweet := Tweet{uid, tweet_id, content, string(time.Now().Format(time.RFC850))}
	t_tweet_content[tweet_id] = tweet
	// 1. 发给自己的tweet里 2. 发给followers的buckets里
	t_posted_by[uid] = append(t_posted_by[uid], tweet_id)

	followers := t_follow[uid]
	for _, follower := range followers {
		t_tweet_bucket[follower] = append(t_tweet_bucket[follower], tweet_id)
	}

	c.JSON(200, gin.H{
		"status": "posted",
	})
}

func deleteTweet(c *gin.Context) {

}

func getTweetsById(c *gin.Context) {

}

func follow(c *gin.Context) {

}

func unfollow(c *gin.Context) {

}

func main() {
	router := gin.Default()
	router.POST("/signup", signup)
	router.POST("/login", login)
	router.POST("/unregister", unregister)
	router.POST("/createNewTweet", createNewTweet)
	router.Run()
}

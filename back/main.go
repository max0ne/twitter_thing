package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type User struct {
	Uid      int
	Username string
	Password string
}

type Tweet struct {
	Tid     int
	Content string
	Time    string
}

var user_registered map[int]User

var t_tweet_content map[int]Tweet

var t_tweet_bucket map[int]int

var t_posted_by map[int]int

var t_follow map[int][]int

var usr_cnt int
var tweet_cnt int

func init() {
	usr_cnt = 0
	tweet_cnt = 0
	user_registered = make(map[int]User)
	t_tweet_content = make(map[int]Tweet)
	t_tweet_bucket = make(map[int]int)
	t_posted_by = make(map[int]int)
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

}

func unregister() {

}

func createNewTweet() {

}

func deleteTweet() {

}

func getTweetsById() {

}

func follow() {

}

func unfollow() {

}

func main() {
	router := gin.Default()
	router.POST("/signup", signup)
	router.GET("/login", login)

	router.Run()
}

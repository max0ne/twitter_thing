package main

import (
	"net/http"

	"github.com/max0ne/twitter_thing/back/db"
	"github.com/max0ne/twitter_thing/back/model"

	"github.com/gin-gonic/gin"
)

type globalTables struct {
	userTable     *db.Table
	tweetTable    *db.Table
	bucketTable   *db.Table
	postedByTable *db.Table
	followTable   *db.Table
}

var tables globalTables

// RESTFul Apis
func signup(c *gin.Context) {
	username := c.PostForm("username")
	password := c.PostForm("password")

	// check whether the user already exists
	oldUser, err := model.GetUser(username, tables.userTable)
	if cerr(c, err) {
		return
	}

	if oldUser != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"status": "user already exists"})
		return
	}

	err = model.SaveUser(model.NewUser(username, password), tables.userTable)
	if cerr(c, err) {
		return
	}
	c.JSON(200, gin.H{
		"status":   "posted",
		"username": username,
	})
}

func login(c *gin.Context) {
	username := c.PostForm("username")
	password := c.PostForm("password")

	user, err := model.GetUser(username, tables.userTable)
	if cerr(c, err) {
		return
	}

	if user != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"status": "user does not exist"})
		return
	}

	if user.Password != password {
		c.JSON(http.StatusUnauthorized, gin.H{"status": "wrong password"})
		return
	}

	c.JSON(200, gin.H{
		"status":   "posted",
		"username": username,
	})
}

func cerr(c *gin.Context, err error) bool {
	if err != nil {
		c.JSON(500, err)
		return true
	}
	return false
}

func unregister(c *gin.Context) {
	username := c.PostForm("username")
	user, err := model.GetUser(username, tables.userTable)
	if cerr(c, err) {
		return
	}
	if user == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"status": "user does not exist"})
		return
	}

	model.DeleteUser(*user, tables.userTable)
	c.JSON(200, gin.H{
		"status":   "posted",
		"username": username,
	})
}

func createNewTweet(c *gin.Context) {
	username := c.PostForm("username")
	content := c.PostForm("content")
	tweet := model.NewTweet(username, content)
	err := model.PublishNewTweet(tweet, tables.followTable, tables.tweetTable, tables.bucketTable, tables.postedByTable)
	if cerr(c, err) {
		return
	}

	c.JSON(200, gin.H{
		"status": "posted",
		"tweet":  tweet,
	})
}

func deleteTweet(c *gin.Context) {

}

func follow(c *gin.Context) {

}

func unfollow(c *gin.Context) {

}

func main() {

	store := db.NewStore()
	tables = globalTables{
		userTable:     store.NewTable("userTable"),
		tweetTable:    store.NewTable("tweetTable"),
		bucketTable:   store.NewTable("bucketTable"),
		postedByTable: store.NewTable("postedByTable"),
		followTable:   store.NewTable("followTable"),
	}

	router := gin.Default()
	router.POST("/signup", signup)
	router.POST("/login", login)
	router.POST("/unregister", unregister)
	router.POST("/createNewTweet", createNewTweet)
	router.Run()
}

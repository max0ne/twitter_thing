package main

import (
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/sessions"
	"github.com/max0ne/twitter_thing/back/db"
	"github.com/max0ne/twitter_thing/back/middleware"
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

func cerr(c *gin.Context, err error) bool {
	if err != nil {
		c.JSON(500, err)
		return true
	}
	return false
}

func sendErr(c *gin.Context, code int, err string) {
	c.JSON(code, gin.H{
		"status": err,
	})
}

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

	sess := sessions.Default(c)
	sess.Set("uname", username)
	if cerr(c, sess.Save()) {
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

	sess := sessions.Default(c)
	sess.Set("uname", username)
	if cerr(c, sess.Save()) {
		return
	}

	c.JSON(200, gin.H{
		"status":   "posted",
		"username": username,
	})
}

func unregister(c *gin.Context) {
	user := middleware.GetUser(c)
	if user == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"status": "user does not exist"})
		return
	}

	model.DeleteUser(*user, tables.userTable)
	c.JSON(200, gin.H{
		"status":   "posted",
		"username": user.Uname,
	})
}

func createNewTweet(c *gin.Context) {
	content := c.PostForm("content")
	tweet := model.NewTweet(middleware.GetUser(c).Uname, content)
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
	user := middleware.GetUser(c)
	tid := c.PostForm("tid")
	tweet, err := model.GetTweet(tid, tables.tweetTable)
	if cerr(c, err) {
		return
	}

	if tweet.Uname != user.Uname {
		sendErr(c, http.StatusUnauthorized, user.Uname+" you are not "+tweet.Uname)
		return
	}

	model.DelTweet(tid, tables.tweetTable)
}

func follow(c *gin.Context) {
	uname := c.PostForm("uname")
	if uname == "" {
		sendErr(c, http.StatusBadRequest, "uname required")
		return
	}

	model.Follow(*middleware.GetUser(c), uname, tables.followTable)
}

func unfollow(c *gin.Context) {
	uname := c.PostForm("uname")
	if uname == "" {
		sendErr(c, http.StatusBadRequest, "uname required")
		return
	}

	if cerr(c, model.UnfollowUser(*middleware.GetUser(c), uname, tables.followTable)) {
		return
	}

	if cerr(c, model.UnfollowUserTweet(*middleware.GetUser(c), uname, tables.bucketTable)) {
		return
	}
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
	router.Use(cors.Default())

	cookieStore := sessions.NewCookieStore([]byte("suer_secret_session_secret"))
	router.Use(sessions.Sessions("defaut_session", cookieStore))
	router.Use(middleware.InjectUser(tables.userTable))

	router.POST("/signup", signup)
	router.POST("/login", login)

	router.POST("/unregister", middleware.RequireLogin, unregister)
	router.POST("/createNewTweet", middleware.RequireLogin, createNewTweet)

	router.Run()
}

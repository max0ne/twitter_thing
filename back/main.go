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

	// vid -> [uid]
	followerTable *db.Table
	// uid -> [vid]
	followingTable *db.Table
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

func sendObj(c *gin.Context, key string, obj interface{}) {
	resp := gin.H{}
	resp[key] = obj
	c.JSON(200, resp)
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
	err := model.PublishNewTweet(tweet, tables.followerTable, tables.tweetTable, tables.bucketTable, tables.postedByTable)
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

	model.Follow(*middleware.GetUser(c), uname, tables.followingTable, tables.followerTable)
}

func unfollow(c *gin.Context) {
	uname := c.PostForm("uname")
	if uname == "" {
		sendErr(c, http.StatusBadRequest, "uname required")
		return
	}

	if cerr(c, model.Unfollow(*middleware.GetUser(c), uname, tables.followingTable, tables.followerTable)) {
		return
	}

	if cerr(c, model.UnfollowUserTweet(*middleware.GetUser(c), uname, tables.bucketTable)) {
		return
	}
}

// get users whom i am following
func getFollowing(c *gin.Context) {
	sendObj(c, "items",
		model.GetUsers(
			model.GetFollowing((*middleware.GetUser(c)).Uname, tables.followingTable), tables.userTable,
		))
}

// get users whom i
func getFollower(c *gin.Context) {
	sendObj(c, "items",
		model.GetUsers(
			model.GetFollowers((*middleware.GetUser(c)).Uname, tables.followerTable), tables.userTable,
		))
}

func getUserTweets(c *gin.Context) {
	uname := c.Param("username")
	if uname == "" {
		sendErr(c, http.StatusBadRequest, "uname required")
		return
	}

	sendObj(c, "items", model.GetUserTweets(uname, tables.tweetTable, tables.postedByTable))
}

func getFeed(c *gin.Context) {
	tweets, err := model.GetUserFeed(middleware.GetUser(c).Uname, tables.tweetTable, tables.bucketTable)
	if cerr(c, err) {
		return
	}

	sendObj(c, "items", tweets)
}

func main() {

	store := db.NewStore()
	tables = globalTables{
		userTable:      store.NewTable("userTable"),
		tweetTable:     store.NewTable("tweetTable"),
		bucketTable:    store.NewTable("bucketTable"),
		postedByTable:  store.NewTable("postedByTable"),
		followerTable:  store.NewTable("followerTable"),
		followingTable: store.NewTable("followingTable"),
	}

	router := gin.Default()
	router.Use(cors.Default())

	cookieStore := sessions.NewCookieStore([]byte("suer_secret_session_secret"))
	router.Use(sessions.Sessions("defaut_session", cookieStore))
	router.Use(middleware.InjectUser(tables.userTable))

	router.POST("/user/signup", signup)
	router.POST("/user/login", login)
	router.POST("/user/unregister", middleware.RequireLogin, unregister)

	router.POST("/user/follow", middleware.RequireLogin, follow)
	router.POST("/user/unfollow", middleware.RequireLogin, unfollow)

	router.GET("/user/following", middleware.RequireLogin, getFollowing)
	router.GET("/user/follower", getFollower)

	router.POST("/tweet/new", middleware.RequireLogin, createNewTweet)
	router.POST("/tweet/del/:tid", middleware.RequireLogin, deleteTweet)

	router.GET("/tweet/user/:username", getUserTweets)
	router.GET("/tweet/feed", middleware.RequireLogin, getFeed)

	router.Run()
}

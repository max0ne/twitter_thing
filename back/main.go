package main

import (
	"fmt"
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/sessions"
	"github.com/max0ne/twitter_thing/back/db"
	"github.com/max0ne/twitter_thing/back/middleware"
	"github.com/max0ne/twitter_thing/back/model"

	"github.com/gin-gonic/gin"
)

// Server - -
type Server struct {
	router *gin.Engine
	store  *db.Store
	tables globalTables
}

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

type loginParam struct {
	Username string `form:"username" json:"username" binding:"required"`
	Password string `form:"password" json:"password" binding:"required"`
}

type newTweetParam struct {
	Content string `form:"content" json:"content" binding:"required"`
}

func cerr(c *gin.Context, err error) bool {
	if err != nil {
		fmt.Println(err)
		c.JSON(500, err)
		return true
	}
	return false
}

func sendErr(c *gin.Context, code int, err string) {
	fmt.Println("err", err)
	c.JSON(code, gin.H{
		"status": err,
	})
}

func sendObj(c *gin.Context, obj interface{}) {
	c.JSON(200, obj)
}

// RESTFul Apis
func (s *Server) signup(c *gin.Context) {
	var param loginParam
	if c.Bind(&param) != nil {
		return
	}
	fmt.Println(param.Username, param.Password)

	// check whether the user already exists
	oldUser, err := model.GetUser(param.Username, s.tables.userTable)
	if cerr(c, err) {
		return
	}

	if oldUser != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"status": "user already exists"})
		return
	}

	err = model.SaveUser(model.NewUser(param.Username, param.Password), s.tables.userTable)
	if cerr(c, err) {
		return
	}

	sess := sessions.Default(c)
	sess.Set("uname", param.Username)
	if cerr(c, sess.Save()) {
		return
	}

	userInDB, err := model.GetUser(param.Username, s.tables.userTable)
	if cerr(c, err) {
		return
	}

	sendObj(c, userInDB)
}

func (s *Server) login(c *gin.Context) {
	var param loginParam
	if c.Bind(&param) != nil {
		return
	}

	user, err := model.GetUserWithPassword(param.Username, s.tables.userTable)
	if cerr(c, err) {
		return
	}

	if user == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"status": "user does not exist"})
		return
	}

	if user.Password != param.Password {
		c.JSON(http.StatusUnauthorized, gin.H{"status": "wrong password"})
		return
	}

	sess := sessions.Default(c)
	sess.Set("uname", param.Username)
	if cerr(c, sess.Save()) {
		return
	}

	user.Password = ""
	sendObj(c, user)
}

func (s *Server) unregister(c *gin.Context) {
	user := middleware.GetUser(c)
	if user == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"status": "user does not exist"})
		return
	}

	model.DeleteUser(*user, s.tables.userTable)
	c.JSON(200, gin.H{
		"status":   "posted",
		"username": user.Uname,
	})
}

func (s *Server) getUser(c *gin.Context) {
	username := c.Param("username")
	user, err := model.GetUser(username, s.tables.userTable)
	if cerr(c, err) {
		return
	}
	if user == nil {
		sendErr(c, http.StatusNotFound, fmt.Sprintf("user %s not found", username))
		return
	}
	sendObj(c, *user)
}

func (s *Server) getCurrentUser(c *gin.Context) {
	sendObj(c, *middleware.GetUser(c))
}

func (s *Server) createNewTweet(c *gin.Context) {
	var param newTweetParam
	if c.Bind(&param) != nil {
		return
	}

	tweet := model.NewTweet(middleware.GetUser(c).Uname, param.Content)
	err := model.PublishNewTweet(&tweet, s.tables.followerTable, s.tables.tweetTable, s.tables.bucketTable, s.tables.postedByTable)
	if cerr(c, err) {
		return
	}

	c.JSON(200, tweet)
}

func (s *Server) deleteTweet(c *gin.Context) {
	user := middleware.GetUser(c)
	tid := c.Param("tid")
	tweet, err := model.GetTweet(tid, s.tables.tweetTable)
	if cerr(c, err) {
		return
	}

	if tweet.Uname != user.Uname {
		sendErr(c, http.StatusUnauthorized, user.Uname+" you are not "+tweet.Uname)
		return
	}

	model.DelTweet(tid, s.tables.tweetTable)
}

func (s *Server) follow(c *gin.Context) {
	uname := c.Param("uname")
	if uname == "" {
		sendErr(c, http.StatusBadRequest, "uname required")
		return
	}

	model.Follow(*middleware.GetUser(c), uname, s.tables.followingTable, s.tables.followerTable)
}

func (s *Server) unfollow(c *gin.Context) {
	uname := c.Param("uname")
	if uname == "" {
		sendErr(c, http.StatusBadRequest, "uname required")
		return
	}

	if cerr(c, model.Unfollow(*middleware.GetUser(c), uname, s.tables.followingTable, s.tables.followerTable)) {
		return
	}

	if cerr(c, model.UnfollowUserTweet(*middleware.GetUser(c), uname, s.tables.bucketTable)) {
		return
	}
}

// get users whom i am following
func (s *Server) getFollowing(c *gin.Context) {
	sendObj(c,
		model.GetUsers(
			model.GetFollowing(c.Param("uname"), s.tables.followingTable), s.tables.userTable,
		))
}

// get users whom i
func (s *Server) getFollower(c *gin.Context) {
	sendObj(c,
		model.GetUsers(
			model.GetFollowers(c.Param("uname"), s.tables.followerTable), s.tables.userTable,
		))
}

func (s *Server) getUserTweets(c *gin.Context) {
	uname := c.Param("username")
	if uname == "" {
		sendErr(c, http.StatusBadRequest, "uname required")
		return
	}

	sendObj(c, model.GetUserTweets(uname, s.tables.tweetTable, s.tables.postedByTable))
}

func (s *Server) getFeed(c *gin.Context) {
	tweets, err := model.GetUserFeed(middleware.GetUser(c).Uname, s.tables.tweetTable, s.tables.bucketTable)
	if cerr(c, err) {
		return
	}

	sendObj(c, tweets)
}

// NewServer - make a server
func NewServer() Server {
	store := db.NewStore()
	tables := globalTables{
		userTable:      store.NewTable("userTable"),
		tweetTable:     store.NewTable("tweetTable"),
		bucketTable:    store.NewTable("bucketTable"),
		postedByTable:  store.NewTable("postedByTable"),
		followerTable:  store.NewTable("followerTable"),
		followingTable: store.NewTable("followingTable"),
	}

	s := Server{
		store:  store,
		tables: tables,
	}
	s.router = s.NewRouter()
	return s
}

// NewRouter make a router
func (s *Server) NewRouter() *gin.Engine {
	router := gin.Default()
	router.Use(cors.Default())

	cookieStore := sessions.NewCookieStore([]byte("suer_secret_session_secret"))
	router.Use(sessions.Sessions("ts", cookieStore))
	router.Use(middleware.InjectUser(s.tables.userTable))

	router.POST("/user/signup", s.signup)
	router.POST("/user/login", s.login)
	router.POST("/user/unregister", middleware.RequireLogin, s.unregister)
	router.GET("/user/get/:username", s.getUser)
	router.GET("/user/me", middleware.RequireLogin, s.getCurrentUser)

	router.POST("/user/follow/:uname", middleware.RequireLogin, s.follow)
	router.POST("/user/unfollow/:uname", middleware.RequireLogin, s.unfollow)

	router.GET("/user/following/:uname", s.getFollowing)
	router.GET("/user/follower/:uname", s.getFollower)

	router.POST("/tweet/new", middleware.RequireLogin, s.createNewTweet)
	router.POST("/tweet/del/:tid", middleware.RequireLogin, s.deleteTweet)

	router.GET("/tweet/user/:username", s.getUserTweets)
	router.GET("/tweet/feed", middleware.RequireLogin, s.getFeed)

	if gin.IsDebugging() {
		router.GET("/db", func(c *gin.Context) {
			sendObj(c, s.store.GetM())
		})
	}
	return router
}

func main() {
	NewServer().router.Run()
}

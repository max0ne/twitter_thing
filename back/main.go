package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/max0ne/twitter_thing/back/config"
	"github.com/max0ne/twitter_thing/back/db"
	"github.com/max0ne/twitter_thing/back/middleware"
	"github.com/max0ne/twitter_thing/back/model"
	"github.com/max0ne/twitter_thing/back/util"

	"github.com/gin-gonic/gin"
)

// Server - -
type Server struct {
	router   *gin.Engine
	dbClient *db.Client
	tables   globalTables
}

type globalTables struct {
	userTable     *db.Table
	miscTable     *db.Table
	tweetTable    *db.Table
	bucketTable   *db.Table
	postedByTable *db.Table

	// vid -> [uid]
	followerTable *db.Table
	// uid -> [vid]
	followingTable *db.Table
}

type loginParam struct {
	Uname    string `form:"uname" json:"uname" binding:"required"`
	Password string `form:"password" json:"password" binding:"required"`
}

type newTweetParam struct {
	Content string `form:"content" json:"content" binding:"required"`
}

type loginResponse struct {
	Uname string `json:"uname"`
	Token string `json:"token"`
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

func sendObjOrErr(c *gin.Context, code int) func(obj interface{}, err error) {
	return func(obj interface{}, err error) {
		if err != nil {
			sendErr(c, code, err.Error())
		} else {
			sendObj(c, obj)
		}
	}
}

// RESTFul Apis
func (s *Server) signup(c *gin.Context) {
	var param loginParam
	if c.Bind(&param) != nil {
		return
	}

	// check whether the user already exists
	oldUser, err := model.GetUser(param.Uname, s.tables.userTable)
	if cerr(c, err) {
		return
	}

	if oldUser != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"status": "user already exists"})
		return
	}

	err = model.SaveUser(model.NewUser(param.Uname, param.Password), s.tables.userTable, s.tables.miscTable)
	if cerr(c, err) {
		return
	}

	token, err := middleware.GenerateJWTToken(param.Uname)
	if cerr(c, err) {
		return
	}

	c.Writer.Header().Set(middleware.TokenHeader, token)
	sendObj(c, loginResponse{
		Uname: param.Uname,
		Token: token,
	})
}

func (s *Server) login(c *gin.Context) {
	var param loginParam
	if c.Bind(&param) != nil {
		return
	}

	user, err := model.GetUserWithPassword(param.Uname, s.tables.userTable)
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

	token, err := middleware.GenerateJWTToken(param.Uname)
	if cerr(c, err) {
		return
	}

	c.Writer.Header().Set(middleware.TokenHeader, token)
	sendObj(c, loginResponse{
		Uname: param.Uname,
		Token: token,
	})
}

func (s *Server) unregister(c *gin.Context) {
	user := middleware.GetUser(c)
	if user == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"status": "user does not exist"})
		return
	}

	model.DeleteUser(*user, s.tables.userTable)
	model.DeleteAllUsersTweet(user.Uname, s.tables.tweetTable, s.tables.postedByTable)
	c.JSON(200, gin.H{
		"status": "posted",
		"uname":  user.Uname,
	})
}

func (s *Server) getUser(c *gin.Context) {
	uname := c.Param("uname")
	user, err := model.GetUser(uname, s.tables.userTable)
	if cerr(c, err) {
		return
	}
	if user == nil {
		sendErr(c, http.StatusNotFound, fmt.Sprintf("user %s not found", uname))
		return
	}
	sendObj(c, *user)
}

func (s *Server) getCurrentUser(c *gin.Context) {
	token, err := middleware.GenerateJWTToken((*middleware.GetUser(c)).Uname)
	if cerr(c, err) {
		return
	}

	c.Writer.Header().Set(middleware.TokenHeader, token)
	sendObj(c, loginResponse{
		Uname: (*middleware.GetUser(c)).Uname,
		Token: token,
	})
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
	if tweet == nil {
		sendErr(c, http.StatusNotFound, "tweet not exist "+tid)
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
	unames, err := model.GetFollowing(c.Param("uname"), s.tables.followingTable)
	if cerr(c, err) {
		return
	}

	sendObj(c,
		model.GetUsers(
			unames, s.tables.userTable,
		))
}

// get users whom i
func (s *Server) getFollower(c *gin.Context) {
	unames, err := model.GetFollowers(c.Param("uname"), s.tables.followerTable)
	if cerr(c, err) {
		return
	}
	sendObj(c,
		model.GetUsers(
			unames, s.tables.userTable,
		))
}

func (s *Server) getUserTweets(c *gin.Context) {
	uname := c.Param("uname")
	if uname == "" {
		sendErr(c, http.StatusBadRequest, "uname required")
		return
	}

	sendObjOrErr(c, http.StatusInternalServerError)(model.GetUserTweets(uname, s.tables.tweetTable, s.tables.postedByTable))
}

func (s *Server) getFeed(c *gin.Context) {
	tweets, err := model.GetUserFeed(middleware.GetUser(c).Uname, s.tables.tweetTable, s.tables.bucketTable)
	if cerr(c, err) {
		return
	}

	sendObj(c, tweets)
}

func (s *Server) getNewRegisterUsers(c *gin.Context) {
	unames, err := model.GetNewRegisteredUserNames(s.tables.miscTable)
	if cerr(c, err) {
		return
	}
	requestUser := middleware.GetUser(c)
	if requestUser != nil {
		unames = util.Remove(unames, requestUser.Uname)
	}
	users := model.GetUsers(unames, s.tables.userTable)
	sendObj(c, users)
}

// NewServer - make a server
func NewServer(config config.Config) Server {
	dbClient, err := db.NewClient(config.DBURL())
	if err != nil {
		log.Fatal(err)
	}
	tables := globalTables{
		miscTable:      dbClient.NewTable("miscTable"),
		userTable:      dbClient.NewTable("userTable"),
		tweetTable:     dbClient.NewTable("tweetTable"),
		bucketTable:    dbClient.NewTable("bucketTable"),
		postedByTable:  dbClient.NewTable("postedByTable"),
		followerTable:  dbClient.NewTable("followerTable"),
		followingTable: dbClient.NewTable("followingTable"),
	}

	s := Server{
		dbClient: dbClient,
		tables:   tables,
	}
	s.router = s.NewRouter()
	return s
}

// NewRouter make a router
func (s *Server) NewRouter() *gin.Engine {
	router := gin.Default()
	router.Use(cors.New(cors.Config{
		AllowMethods:     []string{"GET", "POST", "PUT", "HEAD"},
		AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type", middleware.TokenHeader},
		AllowCredentials: false,
		AllowAllOrigins:  true,
		MaxAge:           12 * time.Hour,
	}))

	router.Use(middleware.InjectUser(s.tables.userTable))

	router.POST("/user/signup", s.signup)
	router.POST("/user/login", s.login)
	router.POST("/user/unregister", middleware.RequireLogin, s.unregister)
	router.GET("/user/get/:uname", s.getUser)
	router.GET("/user/me", middleware.RequireLogin, s.getCurrentUser)
	router.GET("/user/new", s.getNewRegisterUsers)

	router.POST("/user/follow/:uname", middleware.RequireLogin, s.follow)
	router.POST("/user/unfollow/:uname", middleware.RequireLogin, s.unfollow)

	router.GET("/user/following/:uname", s.getFollowing)
	router.GET("/user/follower/:uname", s.getFollower)

	router.POST("/tweet/new", middleware.RequireLogin, s.createNewTweet)
	router.POST("/tweet/del/:tid", middleware.RequireLogin, s.deleteTweet)

	router.GET("/tweet/user/:uname", s.getUserTweets)
	router.GET("/tweet/feed", middleware.RequireLogin, s.getFeed)

	if gin.IsDebugging() {
		router.GET("/db", func(c *gin.Context) {
			sendObjOrErr(c, 500)(s.dbClient.GetM())
		})
	}
	return router
}

func run(config config.Config) error {
	switch config.Role {
	case "api":
		return NewServer(config).router.Run()
	case "db":
		_, err := db.RunServer(config)
		return err
	default:
		return fmt.Errorf("illegal role " + config.Role)
	}
}

func main() {
	config := config.Config{
		Role:       util.GetEnvMust("Role"),
		DBAddr:     util.GetEnvMust("DBAddr"),
		DBPort:     util.GetEnvMust("DBPort"),
		VRPort:     util.GetEnvMust("VRPort"),
		VRPeerURLs: strings.Split(util.GetEnvMust("VRPeerURLs"), ","),
		DBPeerURLs: strings.Split(util.GetEnvMust("DBPeerURLs"), ","),
	}
	err := run(config)
	if err != nil {
		log.Fatal(err)
	}
	for {
		<-time.After(time.Hour)
	}
}

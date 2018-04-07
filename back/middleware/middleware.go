package middleware

import (
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/max0ne/twitter_thing/back/db"
	"github.com/max0ne/twitter_thing/back/model"
)

// InjectUser session -> db -> c.Get("user")
func InjectUser(userTable *db.Table) gin.HandlerFunc {
	return func(c *gin.Context) {
		sess := sessions.Default(c)
		uname := sess.Get("uname")
		if uname == nil {
			return
		}
		username, ok := uname.(string)
		if !ok {
			return
		}

		user, _ := model.GetUser(username, userTable)
		if user != nil {
			c.Set("user", *user)
		}
	}
}

// GetUser get user from context
func GetUser(c *gin.Context) *model.User {
	user, _ := c.Get("user")

	if user == nil {
		return nil
	}
	uu, ok := user.(model.User)
	if !ok {
		return nil
	}

	return &uu
}

// RequireLogin a middleware that sends 401 if not logged in
func RequireLogin(c *gin.Context) {
	user := GetUser(c)
	if user == nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"status": "login required"})
		return
	}
	c.Next()
}

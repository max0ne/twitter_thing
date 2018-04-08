package middleware

import (
	"fmt"
	"net/http"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/max0ne/twitter_thing/back/db"
	"github.com/max0ne/twitter_thing/back/model"
)

const hmacSecret = "hahahhahah"

// TokenHeader - -
const TokenHeader = "X-Twitter-Thing-Token"

// GenerateJWTToken - -
func GenerateJWTToken(uname string) (string, error) {
	return jwt.NewWithClaims(jwt.SigningMethodHS256, &jwt.StandardClaims{
		Id: uname,
	}).SignedString([]byte(hmacSecret))
}

// InjectUser jwt token -> HMAC validate -> db validate -> c.Get("user")
func InjectUser(userTable *db.Table) gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.Request.Header.Get(TokenHeader)
		if tokenString == "" {
			return
		}

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			// Don't forget to validate the alg is what you expect:
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(hmacSecret), nil
		})

		if err != nil {
			fmt.Println(err)
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok || !token.Valid {
			fmt.Println("token not valid")
			return
		}

		uname, ok := claims["jti"].(string)
		if !ok {
			return
		}
		user, err := model.GetUser(uname, userTable)
		if user == nil || err != nil {
			return
		}
		c.Set("user", *user)
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

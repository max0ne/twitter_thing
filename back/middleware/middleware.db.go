package main

import "github.com/gin-gonic/gin"

func DBMiddleware(tables *Tables) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("db")
	}
}

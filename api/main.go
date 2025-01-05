package api

import (
	"github.com/gin-gonic/gin"
)

func StartApi() {
	r := gin.Default()
	r.GET("/", func(c *gin.Context) {
		c.String(200, "hello from gin!")
	})

	// Sub
	r.GET("/sub", handleGetSubApi)

	// Users
	r.GET("/user/:apiToken/:id", handleGetUserApi)

	// Listen on port 8080 by default
	r.Run()
}

package api

import (
	"github.com/gin-gonic/gin"
)

func StartApi() {
	r := gin.Default()
	r.GET("/", func(c *gin.Context) {
		c.String(200, "hello from gin!")
	})
	r.GET("/api/v1/sub", HandleSubApi)

	// Listen on port 8080 by default
	r.Run()
}

package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func Register(router *gin.Engine) {
	apiV1 := router.Group("/api/v1")
	{
		apiV1.GET("/users/:id", func(c *gin.Context) {
			c.JSON(http.StatusAccepted, nil)
		})
	}
}

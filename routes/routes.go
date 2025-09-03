package routes

import (
	"gorm.io/gorm"

	"github.com/gin-gonic/gin"
	"github.com/wolbyte/go_otp/handlers"
)

func Register(router *gin.Engine, db *gorm.DB) {
	apiV1 := router.Group("/api/v1")
	{
		userHandler := handlers.NewUserHandler(db)
		oauthHandler := handlers.NewOAuthHandler(db)

		apiV1.GET("/users/:id", userHandler.GetUser)
		apiV1.POST("/users/oauth", oauthHandler.OAuth)
	}
}

package routes

import (
	"log"
	"strconv"
	"time"

	"gorm.io/gorm"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	_ "github.com/wolbyte/go_otp/docs"
	"github.com/wolbyte/go_otp/utils"

	"github.com/gin-gonic/gin"
	"github.com/wolbyte/go_otp/handlers"
	"github.com/wolbyte/go_otp/middleware"
)

func Register(router *gin.Engine, db *gorm.DB) {
	enableSwagger, err := strconv.ParseBool(utils.GetenvDefault("SERVE_SWAGGER", "true"))

	if err != nil {
		log.Println("Failed to read swagger env, defaulting to true")
		enableSwagger = true
	}

	if enableSwagger && gin.Mode() == gin.ReleaseMode {
		log.Println("WARNING: swagger was enabled in release, refusing to serve.")
		enableSwagger = false
	}

	if enableSwagger {
		log.Println("Serving swagger at localhost:8080/swagger/index.html")
		router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	}

	apiV1 := router.Group("/api/v1")
	{
		userHandler := handlers.NewUserHandler(db)
		oauthHandler := handlers.NewOAuthHandler(db)

		apiV1.GET("/users", userHandler.GetUsers)
		apiV1.GET("/users/:id", userHandler.GetUser)
		apiV1.POST("/users/oauth", middleware.RateLimitOTP(3, 10*time.Minute), oauthHandler.OAuth)

		profileApi := apiV1.Group("/profile")
		profileApi.Use(middleware.AuthRequired())
		{
			profileApi.GET("/info", userHandler.GetProfile)
		}
	}
}

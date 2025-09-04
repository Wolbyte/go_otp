package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/wolbyte/go_otp/utils"
)

// AuthRequired ensures that the request has a valid JWT
func AuthRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			utils.NewHttpError(c, http.StatusUnauthorized, "missing or invalid token")
			c.Abort()
			return
		}

		tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
		claims, err := utils.ParseJWT(tokenStr)
		if err != nil {
			utils.NewHttpError(c, http.StatusUnauthorized, "invalid token or expired token")
			c.Abort()
			return
		}

		// Store user ID in context
		c.Set("user_id", claims.UserID)
		c.Next()
	}
}

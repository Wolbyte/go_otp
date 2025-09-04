package middleware

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/wolbyte/go_otp/handlers"
)

type ClientLimiter struct {
	reqCount uint
	limitEnd time.Time
}

var (
	clients      = make(map[string]*ClientLimiter)
	clientsMutex sync.Mutex
)

func ClientCleanup() {
	for {
		time.Sleep(time.Minute)

		clientsMutex.Lock()
		for id, client := range clients {
			if time.Now().After(client.limitEnd) {
				delete(clients, id)
			}
		}
		clientsMutex.Unlock()
	}
}

func RateLimitOTP(reqLimit uint, limitDuration time.Duration) gin.HandlerFunc {
	go ClientCleanup()

	return func(c *gin.Context) {
		var req handlers.OAuthRequest

		if err := c.ShouldBindBodyWith(&req, binding.JSON); err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": err})
			return
		}

		clientsMutex.Lock()
		client, exists := clients[req.PhoneNumber]
		if !exists {
			clients[req.PhoneNumber] = &ClientLimiter{
				reqCount: 1,
				limitEnd: time.Now().Add(limitDuration),
			}
			client = clients[req.PhoneNumber]
		}

		if client.reqCount > reqLimit {
			clientsMutex.Unlock()
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{"message": "too many requests"})
			return
		}

		client.reqCount += 1
		clientsMutex.Unlock()

		c.Next()
	}
}

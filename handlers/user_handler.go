package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/wolbyte/go_otp/models"
	"gorm.io/gorm"
)

type UserHandler struct {
	DB *gorm.DB
}

type RegisterRequest struct {
	PhoneNumber string `json:"phoneNumber" binding:"required,min=10,max=11"`
}

func NewUserHandler(db *gorm.DB) *UserHandler {
	return &UserHandler{DB: db}
}

func (h *UserHandler) GetUser(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid user id"})
		return
	}

	// Query user from database
	var user models.User
	if err := h.DB.First(&user, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch user"})
		return
	}

	// Return user (exclude password for security)
	c.JSON(http.StatusOK, gin.H{
		"id":          user.ID,
		"phoneNumber": user.PhoneNumber,
		"created":     user.CreatedAt,
	})
}

func (h *UserHandler) GetProfile(c *gin.Context) {
	userID := c.GetUint("user_id")
	c.JSON(http.StatusAccepted, gin.H{"message": "private data", "user_id": userID})
}

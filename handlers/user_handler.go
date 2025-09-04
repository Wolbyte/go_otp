package handlers

import (
	"fmt"
	"math"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/wolbyte/go_otp/models"
	"github.com/wolbyte/go_otp/utils"
	"gorm.io/gorm"
)

type UserHandler struct {
	DB *gorm.DB
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
		"id":            user.ID,
		"phone_number":  user.PhoneNumber,
		"registered_at": user.RegisteredAt,
	})
}

func (h *UserHandler) GetUsers(c *gin.Context) {
	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))

	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	pageSize, err := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	if page < 1 {
		page = 1
	}

	pageSize = utils.Clamp(pageSize, 1, 100)

	dbOffset := (page - 1) * pageSize

	phoneNumber := c.Query("phone_number")
	dateFrom := c.Query("date_from")
	dateTo := c.Query("date_to")
	timezone := c.Query("tz")

	from, to, err := utils.ParseDateRange(dateFrom, dateTo, timezone)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid date"})
		return
	}

	dbQuery := h.DB.Model(&models.User{})

	if phoneNumber != "" {
		dbQuery.Where("phone_number = ?", phoneNumber)
	}

	fmt.Println(from, to)

	if !from.IsZero() && !to.IsZero() {
		dbQuery = dbQuery.Where("registered_at BETWEEN ? AND ?", from, to)
	} else if !from.IsZero() {
		dbQuery = dbQuery.Where("registered_at >= ?", from)
	} else if !to.IsZero() {
		dbQuery = dbQuery.Where("registered_at <= ?", to)
	}

	var users []models.User
	var total int64
	dbQuery.Count(&total)

	result := dbQuery.Offset(dbOffset).Limit(pageSize).Find(&users)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": result.Error.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":          users,
		"page":          page,
		"page_size":     pageSize,
		"total_results": total,
		"total_pages":   int(math.Ceil(float64(total) / float64(pageSize))),
	})
}

func (h *UserHandler) GetProfile(c *gin.Context) {
	userID := c.GetUint("user_id")
	c.JSON(http.StatusAccepted, gin.H{"message": "private data", "user_id": userID})
}

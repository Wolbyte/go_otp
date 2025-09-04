package handlers

import (
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

type GetUsersResponse struct {
	Users       []models.User `json:"users"`
	Page        int           `json:"page" example:"1"`
	PageSize    int           `json:"page_size" example:"10"`
	ResultCount int64         `json:"result_count" example:"30"`
	PageCount   int           `json:"page_count" example:"50"`
}

type GetProfileResponse struct {
	Content string `json:"content" example:"private data!"`
	UserID  uint   `json:"user_id" example:"142"`
}

func NewUserHandler(db *gorm.DB) *UserHandler {
	return &UserHandler{DB: db}
}

// @Summary Get user by ID
// @Description Returns a single user by using it's id as a path argument
// @Tags Users
// @Produce json
// @Param id path int true "ID of the user"
// @Success 200 {object} models.User
// @Failure 404 {object} utils.HTTPError
// @Failure 500 {object} utils.HTTPError
// @Router /users/{id} [get]
func (h *UserHandler) GetUser(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		utils.NewHttpError(c, http.StatusBadRequest, "invalid user id")
		return
	}

	// Query user from database
	var user models.User
	if err := h.DB.First(&user, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			utils.NewHttpError(c, http.StatusNotFound, "user does not exist")
			return
		}
		utils.NewHttpError(c, http.StatusInternalServerError, "failed to fetch user")
		return
	}

	c.JSON(http.StatusOK, user)
}

// @Summary	Get a list of users
// @Description Returns a list of users with pagination and filtering capabilities
// @Tags Users
// @Produce json
// @Param page         query int    false "Sets the current page"
// @Param page_size    query int    false "Elements per page"
// @Param phone_number query string false "Search for user(s) with phone number"
// @Param date_from    query string false "Set the starting registration date (yyyy-mm-dd)" example(2025-09-05)
// @Param date_to      query string false "Set the ending registration date (yyyy-mm-dd)" example(2025-09-06)
// @Param tz           query string false "Parse dates in the desired timezone" example(Asia/Tehran)
// @Success 200 {object} GetUsersResponse
// @Failure 400 {object} utils.HTTPError
// @Failure 500 {object} utils.HTTPError
// @Router /users [get]
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
		utils.NewHttpError(c, http.StatusBadRequest, "invalid date")
		return
	}

	dbQuery := h.DB.Model(&models.User{})

	if phoneNumber != "" {
		dbQuery.Where("phone_number = ?", phoneNumber)
	}

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

	c.JSON(http.StatusOK, GetUsersResponse{
		Users:       users,
		Page:        page,
		PageSize:    pageSize,
		ResultCount: total,
		PageCount:   int(math.Ceil(float64(total) / float64(pageSize))),
	})
}

// @Summary	Get private user data
// @Description Returns a user's id with a message using a JWT token
// @Tags Users
// @Produce json
// @Param Authorization header string true "Bearer Token"
// @Success 200 {object} GetProfileResponse
// @Failure 401 {object} utils.HTTPError
// @Router /profile/info [get]
func (h *UserHandler) GetProfile(c *gin.Context) {
	c.JSON(http.StatusAccepted, GetProfileResponse{
		Content: "private data",
		UserID:  c.GetUint("user_id"),
	})
}

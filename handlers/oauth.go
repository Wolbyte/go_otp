package handlers

import (
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/wolbyte/go_otp/models"
	"github.com/wolbyte/go_otp/utils"
	"gorm.io/gorm"
)

type OAuthHandler struct {
	DB *gorm.DB
}

type OAuthRequest struct {
	PhoneNumber string `json:"phone_number" binding:"required,min=10,max=11" example:"09012345678"`
	OTPCode     string `json:"otp" binding:"max=4" example:""`
}

type OAuthResponse struct {
	Message string `json:"message" example:"success!"`
	Token   string `json:"token" example:"<JWT_TOKEN>"`
}

type OTPData struct {
	code    string
	expiary time.Time
}

func NewOAuthHandler(db *gorm.DB) *OAuthHandler {
	return &OAuthHandler{DB: db}
}

// Use in-memory storage for OTPs to reduce load on the main database
// Redis should be used as a proper implementation
var (
	otpStore = make(map[string]*OTPData)
	otpMutex sync.Mutex
)

// Checks for expired OTPs every minute
func OTPCleanup() {
	for {
		time.Sleep(time.Minute)
		otpMutex.Lock()
		for phone, otp := range otpStore {
			if time.Now().After(otp.expiary) {
				delete(otpStore, phone)
			}
		}
		otpMutex.Unlock()
	}
}

// @Summary	OTP based Register & Login
// @Description Register/Login using phone number and receive an OTP.
// @Description
// @Description Flow:
// @Description 1. Send a phone number in the request body
// @Description 2. Receive an OTP in the console
// @Description 3. Validate the OTP to register the user or login if it already exists (A JWT token will be given in both cases)
// @Description
// @Description Restrictions:
// @Description * Phone number must be unique (for registration)
// @Description * OTP is valid for 2 minutes
// @Description * You can make 3 requests every 10 minutes
// @Tags Users
// @Accept  json
// @Produce json
// @Param request body OAuthRequest true "Request data"
// @Success 202
// @Success 200 {object} OAuthResponse
// @Failure 400 {object} utils.HTTPError
// @Failure 401 {object} utils.HTTPError
// @Failure 429 {object} utils.HTTPError
// @Failure 500 {object} utils.HTTPError
// @Router /users/oauth [post]
func (h *OAuthHandler) OAuth(c *gin.Context) {
	go OTPCleanup()

	var req OAuthRequest

	if err := c.ShouldBindBodyWith(&req, binding.JSON); err != nil {
		utils.NewHttpError(c, http.StatusBadRequest, err.Error())
		return
	}

	validatedNumber, isValidNumber := utils.ValidatePhoneNumber(req.PhoneNumber)
	req.PhoneNumber = validatedNumber

	if !isValidNumber {
		utils.NewHttpError(c, http.StatusBadRequest, "invalid number")
		return
	}

	if req.PhoneNumber != "" && req.OTPCode == "" {
		otpStore[req.PhoneNumber] = &OTPData{code: utils.GenerateOTPCode(), expiary: utils.GenerateExpiary(2)}
		c.JSON(http.StatusAccepted, gin.H{})
		return
	}

	if req.PhoneNumber != "" && req.OTPCode != "" {
		if otpStore[req.PhoneNumber] == nil || otpStore[req.PhoneNumber].code != req.OTPCode {
			utils.NewHttpError(c, http.StatusUnauthorized, "invalid otp")
			return
		}

		otpMutex.Lock()
		delete(otpStore, req.PhoneNumber)
		otpMutex.Unlock()

		var user models.User
		h.DB.Where("phone_number = ?", req.PhoneNumber).First(&user)

		if user.PhoneNumber == "" {
			user.PhoneNumber = req.PhoneNumber

			if err := h.DB.Create(&user).Error; err != nil {
				utils.NewHttpError(c, http.StatusInternalServerError, err.Error())
				return
			}
		}

		token, err := utils.GenerateJWT(user.ID)

		if err != nil {
			utils.NewHttpError(c, http.StatusInternalServerError, fmt.Sprintf("token generation failed: %s", err.Error()))
			return
		}

		c.JSON(http.StatusOK, OAuthResponse{
			Message: "success!",
			Token:   token,
		})
	}
}

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

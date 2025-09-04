package utils

import "github.com/gin-gonic/gin"

type HTTPError struct {
	Code    int    `json:"code" example:"400"`
	Message string `json:"message" example:"page not found"`
}

func NewHttpError(ctx *gin.Context, status int, errorMsg string) {
	err := HTTPError{
		Code:    status,
		Message: errorMsg,
	}

	ctx.JSON(status, err)
}

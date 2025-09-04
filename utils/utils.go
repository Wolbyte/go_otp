package utils

import (
	"fmt"
	"math/rand"
	"strings"
	"time"
)

func GenerateOTPCode() string {
	code := fmt.Sprintf("%04d", rand.Intn(10000))
	fmt.Println("Generated a new OTP:", code)
	return code
}

func ValidatePhoneNumber(number string) (string, bool) {
	isValid := true

	if number == "" {
		return "", false
	}

	number = strings.TrimPrefix(number, "+98")
	number = strings.TrimPrefix(number, "0")

	if strings.Index(number, "9") != 0 {
		isValid = false
	}

	return number, isValid
}

func GenerateExpiary(minutes time.Duration) time.Time {
	return time.Now().Add(minutes * time.Minute)
}

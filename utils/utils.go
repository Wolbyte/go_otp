package utils

import (
	"fmt"
	"math/rand"
	"time"
)

func GenerateOTPCode() string {
	code := fmt.Sprintf("%04d", rand.Intn(10000))
	fmt.Println("Generated a new OTP:", code)
	return code
}

func ValidatePhoneNumber(number string) (string, bool) {
	isValid := true

	if number[0:3] == "+98" {
		number = number[3:]
	}

	if number[0] == '0' {
		number = number[1:]
	}

	if number[0] != '9' {
		isValid = false
	}

	return number, isValid
}

func GenerateExpiary(minutes time.Duration) time.Time {
	return time.Now().Add(minutes * time.Minute)
}

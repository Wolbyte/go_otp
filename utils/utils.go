package utils

import (
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"
)

func GenerateOTPCode() string {
	code := fmt.Sprintf("%04d", rand.Intn(10000))
	fmt.Println("Generated a new OTP:", code)
	return code
}

func GetenvDefault(variable string, defaultValue string) string {
	result := os.Getenv(variable)

	if result == "" {
		result = defaultValue
	}

	return result
}

func Clamp(n int, min int, max int) int {
	if n > max {
		n = max
	} else if n < min {
		n = min
	}

	return n
}

func ParseDateRange(dateFrom, dateTo, tz string) (time.Time, time.Time, error) {
	loc, err := time.LoadLocation(tz)
	if err != nil {
		loc = time.UTC
	}

	var from, to time.Time

	if dateFrom != "" {
		t, err := time.ParseInLocation("2006-01-02", dateFrom, loc)
		if err != nil {
			return from, to, err
		}
		from = t.In(time.UTC) // convert to UTC
	}

	if dateTo != "" {
		t, err := time.ParseInLocation("2006-01-02", dateTo, loc)
		if err != nil {
			return from, to, err
		}
		// end of day in user timezone, then convert to UTC
		to = t.AddDate(0, 0, 1).Add(-time.Nanosecond).In(time.UTC)
	}

	return from, to, nil
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

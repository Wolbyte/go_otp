package utils

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JWTClaims struct {
	UserID uint `json:"user_id"`
	jwt.RegisteredClaims
}

func GetJWTSecret() ([]byte, error) {
	secret := []byte(GetenvDefault("JWT_SECRET", ""))
	if len(secret) == 0 {
		return []byte{}, errors.New("JWT_SECRET was empty")
	}
	return secret, nil
}

func GenerateJWT(userID uint) (string, error) {
	secret, err := GetJWTSecret()
	if err != nil {
		return "", err
	}

	expiration := time.Now().Add(24 * time.Hour)

	claims := &JWTClaims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiration),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(secret)
}

// Validates a token and returns claims
func ParseJWT(tokenStr string) (*JWTClaims, error) {
	secret, err := GetJWTSecret()
	if err != nil {
		return nil, err
	}

	token, err := jwt.ParseWithClaims(tokenStr, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		return secret, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*JWTClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, err
}

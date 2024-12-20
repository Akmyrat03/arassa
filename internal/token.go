package internal

import (
	"errors"

	"github.com/dgrijalva/jwt-go"
)

var (
	signingKey = []byte("###%4544566656")
)

type tokenClaims struct {
	jwt.StandardClaims

	Admin_id int    `json:"admin_id"`
	Username string `json:"username"`
}

// Validate token
func ValidateToken(tokenString string) (string, error) {
	claims := &tokenClaims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {

		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return signingKey, nil
	})

	if err != nil || !token.Valid {
		return "", errors.New("invalid token")
	}

	return claims.Username, nil
}

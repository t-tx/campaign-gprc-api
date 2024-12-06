package pkg

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

type campaignClaims struct {
	Username string            `json:"username"`
	Data     map[string]string `json:"data"`
	jwt.RegisteredClaims
}

func GenerateJWT(username string, data map[string]string, jwtSecret []byte) (string, error) {
	var c = &campaignClaims{
		Username: username,
		Data:     data,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, c)

	return token.SignedString(jwtSecret)
}

func VerifyJWT(tokenString string, jwtSecret []byte) (*campaignClaims, error) {
	// Parse and validate the token
	token, err := jwt.ParseWithClaims(tokenString, &campaignClaims{}, func(token *jwt.Token) (interface{}, error) {
		// Return the secret key for validation
		return jwtSecret, nil
	})

	if err != nil {
		return nil, err
	}

	// Extract the claims from the token
	if claims, ok := token.Claims.(*campaignClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, fmt.Errorf("invalid token")
}

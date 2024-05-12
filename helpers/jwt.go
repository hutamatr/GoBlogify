package helpers

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func AccessTokenDuration(appEnv string) time.Duration {
	var accessTokenExpired time.Duration
	if appEnv == "production" {
		accessTokenExpired = 5 * time.Minute
	} else {
		accessTokenExpired = 2 * time.Hour
	}
	return accessTokenExpired
}

func GenerateToken(userId int, expired time.Duration, tokenSecret string) (string, error) {
	tokenBuilder := jwt.NewWithClaims(jwt.SigningMethodHS256,
		jwt.MapClaims{
			"exp": time.Now().Add(expired).Unix(),
			"iat": time.Now().Unix(),
			"sub": userId,
		})

	tokenString, err := tokenBuilder.SignedString([]byte(tokenSecret))

	return tokenString, err
}

func VerifyToken(tokenString string, tokenSecret []byte) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return tokenSecret, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token")
}
